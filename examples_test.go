package config_test

import (
	"time"

	"github.com/gernoteger/log15-config"
	"gopkg.in/yaml.v2"
)

func getMapFromConfiguration(config string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(config), &configMap)
	if err != nil {
		return nil, err
	}
	return configMap, err
}

func Example() {
	var exampleConfiguration = `
  # default for all handlers
  level: INFO

  # optionally buffer up to bufsize messages; if omitted (or 0) we don't buffer
  bufsize: 100
  extra:
      mark: test
      user: alice

  handlers:
    - kind: stdout
      format: terminal

    - kind: stderr
      format: json
      level: warn	# don't show

    - kind: stdout
      format: logfmt
      level: debug

    - kind: syslog      # syslog is only available on linux
      tag: testing
      facility: local6

    # 2 ways to configure net
    - kind: net
      url: udp://localhost:4242
      format: json
      level: debug
    - kind: net
      url: tcp://localhost:4242
      format: json
      level: debug

    - kind: multi
      handlers:
        - kind: stdout
          format: terminal
        - kind: stderr
          format: json
        - kind: stdout
          format: logfmt

    - kind: filter  # MatchFilterHandler
      key: matcher
      value: foo
      handler:
        kind: stdout
        format: json

    - kind: failover
      handlers:
        - kind: stdout
          format: terminal
        - kind: stderr
          format: json
        - kind: stdout
          format: logfmt

    # a buffered handler encloses other. Instead of using this preferably use the bufsize parameter above to enclose
    # the whole tree into a buffered handler instead.
    - kind: buffer
      level: debug # w/o this, the nested handler(s) won't be activated!!
      bufsize: 100
      handler:
        kind: net
        url: tcp://localhost:4242
        format: json
`
	//hooks.Register(HandlerConfigType, "failover", NewFailoverConfig)

	configMap, err := getMapFromConfiguration(exampleConfiguration)
	if err != nil {
		panic(err)
	}

	log, err := config.Logger(configMap)
	if err != nil {
		panic(err)
	}

	log.Info("Hello, world!")
	log.Debug("user1", "user", "bob")

	l1 := log.New("user", "carol") // issue in log15! won't override, but use both!
	l1.Debug("about user")

	time.Sleep(100 * time.Millisecond) // need this to finish all async log messages. Bufferedhandler doesn't expose a means to see if the channel is closed...

	// disabling output below for tests by immediately prepending this line since dates will never be right. to execute, just insert blank line after this ons.

	// Output:
	// INFO[11-30|11:37:20] Hello, world!                            mark=test user=alice
	// t=2016-11-30T11:37:20+0100 lvl=info msg="Hello, world!" mark=test user=alice
	// t=2016-11-30T11:37:20+0100 lvl=dbug msg=user1 mark=test user=alice user=bob
	// t=2016-11-30T11:37:20+0100 lvl=dbug msg="about user" mark=test user=alice user=carol

}
