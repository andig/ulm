package cmd

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/evcc-io/evcc/server"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/test"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// chargerCmd represents the charger command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create an EVCC configuration",
	Run:   runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) {
	util.LogLevel(viper.GetString("log"), viper.GetStringMapString("levels"))
	log.INFO.Printf("evcc %s (%s)", server.Version, server.Commit)

	fmt.Println()
	fmt.Println("The next steps will guide throught the creation of a EVCC configuration file.")
	fmt.Println("Please be aware that this process does not cover all possible scenarios.")
	fmt.Println("You can stop the process by pressing ctrl-c.")
	fmt.Println()
	fmt.Println("Let's start:")

	fmt.Println("1. Configure your charger")
	_ = processCharger()
}

type Config struct {
	Chargers []map[string]interface{} `yaml:"chargers"`
}

// let the user select a charger
func processCharger() test.ConfigTemplate {
	var finished bool
	var chargerConfiguration test.ConfigTemplate

	for ok := true; ok; ok = finished {
		var configuration Config

		fmt.Println()
		configItem := selectItem("charger")
		chargerConfiguration = processConfig(configItem, defaultChargerName)

		configuration.Chargers = append(configuration.Chargers, chargerConfiguration.Config)

		yaml, err := yaml.Marshal(configuration)
		if err != nil {
			log.FATAL.Fatalf("Invalid YAML: %s", err)
		}

		// check if we need to setup an EEBUS hems

		fmt.Println()
		fmt.Println("Testing configuration...")
		fmt.Println()

		err = testChargerConfig(yaml)
		if err != nil {
			log.FATAL.Fatalf("Invalid charger configuration: %s", err)
		}

		// Do we see proper values?
		fmt.Println()
		if askYesNo("Does the test data above show proper values?") {
			fmt.Printf("\n2. Final configuration:\n\n%s\n", string(yaml))
			finished = true
		} else {
			fmt.Println()
			if !askYesNo("Do you want to try it again?") {
				finished = true
			}
		}
	}

	return chargerConfiguration
}

// return EVCC configuration items of a given type
func fetchElements(typ string) []test.ConfigTemplate {
	var items []test.ConfigTemplate

	for _, tmpl := range test.ConfigTemplates(typ) {
		items = append(items, tmpl)
	}

	return items
}

// PromptUI: select item from list
func selectItem(typ string) test.ConfigTemplate {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "-> {{ .Name }}",
		Inactive: "   {{ .Name }}",
		Selected: fmt.Sprintf("%s: {{ .Name }}", typ),
	}

	items := fetchElements(typ)
	prompt := promptui.Select{
		Label:     fmt.Sprintf("Choose %s", typ),
		Items:     items,
		Templates: templates,
		Size:      10,
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.FATAL.Fatal(err)
	}

	return items[index]
}

func askYesNo(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		// PromptUI returns an error if the response is not "y"
		return false
	}

	return strings.ToLower(result) == "yes"
}

// PromputUI: ask for input
func askValue(label string, defaultValue interface{}) interface{} {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	validate := func(input string) error {
		return nil
	}
	var defValue string
	switch v := defaultValue.(type) {
	case nil:
		defValue = ""
	case string:
		defValue = v
	case int:
		defValue = strconv.Itoa(v)
		validate = func(input string) error {
			_, err := strconv.ParseInt(input, 10, 64)
			return err
		}
	default:
		log.FATAL.Fatalf("unsupported type: %s", defaultValue)
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Default:   defValue,
		Validate:  validate,
		AllowEdit: true,
	}

	result, err := prompt.Run()
	if err != nil {
		log.FATAL.Fatal(err)
	}

	var returnValue interface{}
	switch defaultValue.(type) {
	case nil:
		returnValue = result
	case string:
		returnValue = result
	case int:
		returnValue, err = strconv.Atoi(result)
		if err != nil {
			log.FATAL.Fatal("entered invalid int value")
		}
	default:
		log.FATAL.Fatalf("unsupported type: %s", defaultValue)
	}
	return returnValue
}

// Process an EVCC configuration item
func processConfig(configItem test.ConfigTemplate, defaultName string) test.ConfigTemplate {
	// check for parameters the user has to provide
	var conf map[string]interface{}
	if err := yaml.Unmarshal([]byte(configItem.Sample), &conf); err != nil {
		// silently ignore errors here
		log.WARN.Printf("unable to parse sample: %s", err)
	}

	parsed := test.ConfigTemplate{
		Template: configItem.Template,
		Config:   conf,
	}

	var enteredURI string

	if len(parsed.Config) > 0 {
		fmt.Println()
		fmt.Println("Enter the configuration values:")
		for param, value := range parsed.Config {
			var prompt string
			valueType := reflect.ValueOf(value)

			switch param {
			case "user":
				prompt = "Username"
			case "password":
				prompt = "Password"
			case "meter":
				// e.g. Discovery meter
				prompt = "Identifier"
			case "device":
				// Serial modbus devices
				prompt = "Serial port"
			case "baudrate":
				// Serial modbus devices
				prompt = "Serial baudrate"
			case "ain":
				// Fritzbox Dect devices
				prompt = "AIN (printed on the device)"
			case "token":
				// Tokens, e.g. go-e Charger
				prompt = "Token"
			case "mac":
				// MAC Address, e.g. NRGKick Charger
				prompt = "MAC address"
			case "pin":
				// PIN Number, e.g. NRGKick Charger
				prompt = "PIN number"
			case "charger":
				// Charger Identifier, e.g. Easee
				prompt = "Charger identifier"
			default:
				if valueType.Kind() == reflect.String {
					// check if the value constains a default URI 192.0.2.2 string
					if strings.Contains(value.(string), defaultURI) {
						if len(enteredURI) == 0 {
							enteredURI = askValue("Address:", defaultURI).(string)
						}

						value = strings.Replace(value.(string), defaultURI, enteredURI, -1)
					}
				}
			}

			if len(prompt) > 0 {
				value = askValue(prompt, value)
			}

			parsed.Config[param] = value
		}
	}

	parsed.Config["type"] = configItem.Template.Type

	fmt.Println()
	fmt.Println("Provide a name for this device:")
	parsed.Config["name"] = askValue("Name", defaultName).(string)

	return parsed
}

// return a usable EVCC configuration
func readConfiguration(configuration []byte) (conf config, err error) {
	if err := viper.ReadConfig(bytes.NewBuffer(configuration)); err != nil {
		return conf, err
	}

	if err := viper.UnmarshalExact(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

// test device configuration
func testChargerConfig(configuration []byte) error {
	var conf config

	conf, err := readConfiguration(configuration)
	if err != nil {
		return err
	}

	if err := cp.configureChargers(conf); err != nil {
		return err
	}

	chargers := cp.chargers

	d := dumper{len: 1}
	for name, v := range chargers {
		d.DumpWithHeader(name, v)
	}

	return nil
}
