package server

import (
	"context"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log"
	"go-blog/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var log = logging.Logger("server")

// Server
/**
 * Interface for server
 *
 * Run(): Run the server. The caller would be blocked unless errors raised.
 * Note that Run() would not start the server.

 * Start(): Send start cmd to main routine. Should be called in a separate go routine.

 * ShutDown(): Send shutdown cmd to main routine. Server would be safely shutdown.
 * Should be called in a separate go routine.

 * Restart(): Send restart cmd to main routine. Server would be restarted with config refreshed.
 * Should be called in a separate go routine.
 */
type Server interface {
	Run() error
	Start()
	ShutDown()
	Restart()
}

type serverCmd uint8
const cmdServerStart serverCmd = 0
const cmdServerShutdown serverCmd = 1
const cmdServerRestart serverCmd = 2

type ginServer struct {
	router	*gin.Engine
	server  *http.Server	// Used to control the lifecycle of server.
	cfg 	config.RunningConfig
	prefix  string			// Used to build url. It depends on running config when initializing.
	isRunning bool
	cmdCh	chan serverCmd
	errCh	chan error
	wg		sync.WaitGroup
	ctx 	context.Context
}

// NewGinServer
// TODO: Use config to control the behavior of engine.
func NewGinServer(cfg config.Config) Server {
	runCfg := cfg.RunningConfig()
	res := &ginServer{
		cfg:    runCfg,
		isRunning: false,
		cmdCh: make(chan serverCmd),
		errCh: make(chan error),
		ctx: context.Background(),
	}
	res.initRouter()
	res.server = &http.Server{
		Addr:    ":" + strconv.Itoa(int(runCfg.Port())),
		Handler: res.router,
	}

	if runCfg.HostOnly() {
		res.prefix = "http://" + runCfg.Host()
	} else {
		res.prefix = "http://" + runCfg.Host() + ":" + strconv.Itoa(int(runCfg.Port()))
	}

	return res
}

func (s *ginServer) reset(cfg config.Config) {
	s.cfg = cfg.RunningConfig()
	s.initRouter()
	s.server = &http.Server{
		Addr: ":" + strconv.Itoa(int(s.cfg.Port())),
		Handler: s.router,
	}
	if s.cfg.HostOnly() {
		s.prefix = "http://" + s.cfg.Host()
	} else {
		s.prefix = "http://" + s.cfg.Host() + ":" + strconv.Itoa(int(s.cfg.Port()))
	}
}

func (s *ginServer) Run() error {
	var err error
	defer func() {
		s.wg.Wait()
		if err != nil {
			log.Error("Error when run server: ", err)
		}
		log.Debug("Server end.")
	}()

	var sCmd serverCmd
	quitCh := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sCmd = <- s.cmdCh:
			switch sCmd {
			case cmdServerStart:
				log.Debug("Server start.")
				if s.isRunning {
					log.Warn("Start when server is already started.")
				} else {
					s.isRunning = true
					go s.startServer()
				}
			case cmdServerRestart:
				log.Debug("Server restart.")
				if s.isRunning {
					err = s.server.Shutdown(s.ctx)
					if err != nil {
						return err
					}
					s.isRunning = false
				}
				// Reopen config
				log.Debug("Reconfigure server.")
				var newCfg config.Config
				newCfg, err = config.OpenFileConfig()
				if err != nil {
					return err
				}
				s.reset(newCfg)
				s.isRunning = true
				go s.startServer()
			case cmdServerShutdown:
				if s.isRunning {
					log.Debug("Server Restart ...")
					err = s.server.Shutdown(s.ctx)
					if err != nil {
						return err
					}
					s.isRunning = false
				} else {
					log.Warn("Server shutdown while it is not running.")
				}
			default:
				log.Error("Unknown server command ", sCmd)
			}

		case err = <- s.errCh:
			return err
		case <- quitCh:
			err = s.server.Shutdown(s.ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (s *ginServer) startServer() {
	var err error
	defer func() {
		log.Debug("Server end with error: ", err)
		s.wg.Done()
	}()
	s.wg.Add(1)
	//err = s.router.Run(":" + strconv.Itoa(int(s.cfg.Port())))
	err = s.server.ListenAndServe()
	if err != nil {
		s.isRunning = false
		if err != http.ErrServerClosed {	// ServerClosed would not end the program
			s.errCh <- err
		}
	}
}

func (s *ginServer) Start() {
	s.cmdCh <- cmdServerStart
}

func (s *ginServer) ShutDown() {
	s.cmdCh <- cmdServerShutdown
}

func (s *ginServer) Restart() {
	s.cmdCh <- cmdServerRestart
}

func (s *ginServer) buildUrl(relativePath string) string {
	return s.prefix + relativePath
}