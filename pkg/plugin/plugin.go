package plugin

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type commonFlags struct {
	allNamespaces      bool                  // should we search all namespaces
	container          string                // name of the container to search for
	filterList         map[string]matchValue // used to filter out rows form the table during Print function
	labels             string                // k8s pod labels
	showInitContainers bool                  // currently only for mem and cpu sub commands, placed here incase its needed in the future for others
	showOddities       bool                  // this isnt really common but it does show up across 3+ commands and im lazy
	showNamespaceName  bool                  // shows the namespace name of each pod
	showNodeName       bool                  // do we need to show the node name in the output
	showTreeView       bool                  // show the table in a tree like view
	showNodeTree       bool                  // show the tree rooted at the node level, forces showTreeView to true
	showContainerType  bool                  // show container type column
	byteSize           string                // sets the bytes conversion for the output size
	outputAs           string                // how to output the table, currently only accepts json
	sortList           []string              // column names to sort on when table.Print() is called
	matchSpecList      map[string]matchValue // filter pods based on matches to the v1.Pods.Spec fields
	calcMatchOnly      bool                  // should we calculate up only the rows that match
	inputFilename      string                // filename to read pod information from, rather than the k8s api
	labelNodeName      string
	labelPodName       string
	annotationPodName  string
	showColumnByName   string // list of column names to show, overrides other hidden columns
	outputAsColour     int    // which coloring type do we use when displaying columns
	useTheseColours    [][2]int
}

const (
	COLOUR_NONE      = 0
	COLOUR_ERRORS    = 1
	COLOUR_COLUMNS   = 2
	COLOUR_MIX       = 3
	COLOUR_CUSTOM    = 4
	COLOUR_CUSTOMMIX = 5
)

func InitSubCommands(rootCmd *cobra.Command) {
	var includeInitShort string = "include init container(s) in the output, by default init containers are hidden"
	var odditiesShort string = "show only the outlier rows that dont fall within the computed range"
	var sizeShort string = "allows conversion to the selected size rather then the default megabyte output"
	var treeShort string = "Display tree like view instead of the standard list"
	var nodetreeShort string = "Displays the tree with the nodes as the root"
	var showIPShort string = "Show the pods IP address column"
	// var treeShort string = "Display tree like view instead of the standard list"

	log := logger{location: "InitSubCommands"}
	log.Debug("Start")

	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)
	rootCmd.SetHelpTemplate(helpTemplate)

	// capabilities
	var cmdCapabilities = &cobra.Command{
		Use:     "capabilities",
		Short:   capabilitiesShort,
		Long:    fmt.Sprintf("%s\n\n%s", capabilitiesShort, capabilitiesDescription),
		Example: fmt.Sprintf(capabilitiesExample, rootCmd.CommandPath()),
		Aliases: []string{"cap"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Capabilities(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdCapabilities.Flags())
	cmdCapabilities.Flags().BoolP("tree", "t", false, treeShort)
	cmdCapabilities.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdCapabilities)
	rootCmd.AddCommand(cmdCapabilities)

	// commands
	var cmdCommands = &cobra.Command{
		Use:     "command",
		Short:   commandsShort,
		Long:    fmt.Sprintf("%s\n\n%s", commandsShort, commandsDescription),
		Example: fmt.Sprintf(commandsExample, rootCmd.CommandPath()),
		Aliases: []string{"cmd", "exec", "args"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Commands(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdCommands.Flags())
	cmdCommands.Flags().BoolP("tree", "t", false, treeShort)
	cmdCommands.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdCommands)
	rootCmd.AddCommand(cmdCommands)

	// cpu
	var cmdCPU = &cobra.Command{
		Use:     "cpu",
		Short:   resourceShort("cpu"),
		Long:    fmt.Sprintf("%s\n\n%s", resourceShort("cpu"), resourceDescription("cpu")),
		Example: fmt.Sprintf(resourceExample("cpu"), rootCmd.CommandPath()),
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Resources(cmd, KubernetesConfigFlags, args, "cpu"); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdCPU.Flags())
	cmdCPU.Flags().BoolP("include-init", "i", false, includeInitShort)
	cmdCPU.Flags().BoolP("oddities", "", false, odditiesShort)
	cmdCPU.Flags().BoolP("raw", "r", false, "show raw values")
	cmdCPU.Flags().BoolP("tree", "t", false, treeShort)
	cmdCPU.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdCPU)
	rootCmd.AddCommand(cmdCPU)

	// environment
	var cmdEnvironment = &cobra.Command{
		Use:     "environment",
		Short:   environmentShort,
		Long:    fmt.Sprintf("%s\n\n%s", environmentShort, environmentDescription),
		Example: fmt.Sprintf(environmentExample, rootCmd.CommandPath()),
		Aliases: []string{"env", "vars"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Environment(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdEnvironment.Flags())
	cmdEnvironment.Flags().BoolP("translate", "", false, "read the configmap show its values")
	cmdEnvironment.Flags().BoolP("tree", "t", false, treeShort)
	cmdEnvironment.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdEnvironment)
	rootCmd.AddCommand(cmdEnvironment)

	// ip
	var cmdIP = &cobra.Command{
		Use:     "ip",
		Short:   ipShort,
		Long:    fmt.Sprintf("%s\n\n%s", ipShort, ipDescription),
		Example: fmt.Sprintf(ipExample, rootCmd.CommandPath()),
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Ports(cmd, KubernetesConfigFlags, args, true); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdIP.Flags())
	addCommonFlags(cmdIP)
	rootCmd.AddCommand(cmdIP)

	// image
	var cmdImage = &cobra.Command{
		Use:     "image",
		Short:   imageShort,
		Long:    fmt.Sprintf("%s\n\n%s", imageShort, imageDescription),
		Example: fmt.Sprintf(imageExample, rootCmd.CommandPath()),
		Aliases: []string{"im"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Image(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdImage.Flags())
	cmdImage.Flags().BoolP("id", "", false, "Show running containers id")
	cmdImage.Flags().BoolP("tree", "t", false, treeShort)
	cmdImage.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdImage)
	rootCmd.AddCommand(cmdImage)

	// lifecycle
	var cmdLifecycle = &cobra.Command{
		Use:     "lifecycle",
		Short:   lifecycleShort,
		Long:    fmt.Sprintf("%s\n\n%s", lifecycleShort, lifecycleDescription),
		Example: fmt.Sprintf(lifecycleExample, rootCmd.CommandPath()),
		Aliases: []string{"im"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Lifecycle(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdLifecycle.Flags())
	cmdLifecycle.Flags().BoolP("tree", "t", false, treeShort)
	cmdLifecycle.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdLifecycle)
	rootCmd.AddCommand(cmdLifecycle)

	// memory
	var cmdMemory = &cobra.Command{
		Use:     "memory",
		Short:   resourceShort("memory"),
		Long:    fmt.Sprintf("%s\n\n%s", resourceShort("memory"), resourceDescription("memory")),
		Example: fmt.Sprintf(resourceExample("memory"), rootCmd.CommandPath()),
		Aliases: []string{"mem"},
		// SuggestFor: []string{""},
		// Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Resources(cmd, KubernetesConfigFlags, args, "memory"); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdMemory.Flags())
	cmdMemory.Flags().BoolP("include-init", "i", false, includeInitShort)
	cmdMemory.Flags().BoolP("oddities", "", false, odditiesShort)
	cmdMemory.Flags().BoolP("raw", "r", false, "show raw values")
	cmdMemory.Flags().String("size", "Mi", sizeShort)
	cmdMemory.Flags().BoolP("tree", "t", false, treeShort)
	cmdMemory.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdMemory)
	rootCmd.AddCommand(cmdMemory)

	// ports
	var cmdPorts = &cobra.Command{
		Use:     "ports",
		Short:   portsShort,
		Long:    fmt.Sprintf("%s\n\n%s", portsShort, portsDescription),
		Example: fmt.Sprintf(portsExample, rootCmd.CommandPath()),
		Aliases: []string{"port", "po"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Ports(cmd, KubernetesConfigFlags, args, false); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdPorts.Flags())
	cmdPorts.Flags().BoolP("tree", "t", false, treeShort)
	cmdPorts.Flags().BoolP("node-tree", "", false, nodetreeShort)
	cmdPorts.Flags().BoolP("show-ip", "", false, showIPShort)
	addCommonFlags(cmdPorts)
	rootCmd.AddCommand(cmdPorts)

	// probes
	var cmdProbes = &cobra.Command{
		Use:     "probes",
		Short:   probesShort,
		Long:    fmt.Sprintf("%s\n\n%s", probesShort, probesDescription),
		Example: fmt.Sprintf(probesExample, rootCmd.CommandPath()),
		Aliases: []string{"probe"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Probes(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdProbes.Flags())
	cmdProbes.Flags().BoolP("tree", "t", false, treeShort)
	cmdProbes.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdProbes)
	rootCmd.AddCommand(cmdProbes)

	// restarts
	var cmdRestart = &cobra.Command{
		Use:     "restarts",
		Short:   restartsShort,
		Long:    fmt.Sprintf("%s\n\n%s", restartsShort, restartsDescription),
		Example: fmt.Sprintf(restartsExample, rootCmd.CommandPath()),
		Aliases: []string{"restart"},
		// SuggestFor: []string{""},
		// Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Restarts(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdRestart.Flags())
	cmdRestart.Flags().BoolP("oddities", "", false, odditiesShort)
	cmdRestart.Flags().BoolP("tree", "t", false, treeShort)
	cmdRestart.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdRestart)
	rootCmd.AddCommand(cmdRestart)

	// security
	var cmdSecurity = &cobra.Command{
		Use:     "security",
		Short:   securityShort,
		Long:    fmt.Sprintf("%s\n\n%s", securityShort, securityDescription),
		Example: fmt.Sprintf(securityExample, rootCmd.CommandPath()),
		Aliases: []string{"sec"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Security(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdSecurity.Flags())
	cmdSecurity.Flags().BoolP("selinux", "", false, "show the SELinux context thats applied to the containers")
	cmdSecurity.Flags().BoolP("tree", "t", false, treeShort)
	cmdSecurity.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdSecurity)
	rootCmd.AddCommand(cmdSecurity)

	// status
	var cmdStatus = &cobra.Command{
		Use:     "status",
		Short:   statusShort,
		Long:    fmt.Sprintf("%s\n\n%s", statusShort, statusDescription),
		Example: fmt.Sprintf(statusExample, rootCmd.CommandPath()),
		Aliases: []string{"st"},
		// SuggestFor: []string{""},
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Status(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdStatus.Flags())
	cmdStatus.Flags().BoolP("details", "d", false, `Display the timestamp instead of age along with the message column`)
	cmdStatus.Flags().BoolP("oddities", "", false, odditiesShort)
	cmdStatus.Flags().BoolP("previous", "p", false, "Show previous state")
	cmdStatus.Flags().BoolP("id", "", false, "Show running containers id")
	cmdStatus.Flags().BoolP("tree", "t", false, treeShort)
	cmdStatus.Flags().BoolP("node-tree", "", false, nodetreeShort)
	// TODO: check if I can add labels for service/replicaset/configmap etc.
	addCommonFlags(cmdStatus)
	rootCmd.AddCommand(cmdStatus)

	// version
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: versionsShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			Version(cmd, KubernetesConfigFlags, args)
			return nil
		},
	}
	rootCmd.AddCommand(cmdVersion)

	// volumes
	var cmdVolume = &cobra.Command{
		Use:     "volumes",
		Short:   volumesShort,
		Long:    fmt.Sprintf("%s\n\n%s", volumesShort, volumesDescription),
		Example: fmt.Sprintf(volumesExample, rootCmd.CommandPath()),
		Aliases: []string{"volume", "vol"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Volumes(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdVolume.Flags())
	cmdVolume.Flags().BoolP("device", "d", false, "show raw block device mappings within a container")
	cmdVolume.Flags().BoolP("tree", "t", false, treeShort)
	cmdVolume.Flags().BoolP("node-tree", "", false, nodetreeShort)
	addCommonFlags(cmdVolume)
	rootCmd.AddCommand(cmdVolume)

}

// adds common flags to the passed command
func addCommonFlags(cmdObj *cobra.Command) {
	cmdObj.Flags().BoolP("all-namespaces", "A", false, "list containers form pods in all namespaces")
	cmdObj.Flags().StringP("selector", "l", "", `Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2`)
	cmdObj.Flags().StringP("container", "c", "", `Container name. If omitted show all containers in the pod`)
	cmdObj.Flags().StringP("sort", "", "", `Sort by column`)
	cmdObj.Flags().StringP("output", "o", "", `Output format, currently csv, list, json and yaml are supported`)
	cmdObj.Flags().StringP("match", "m", "", `Filters out results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != `)
	cmdObj.Flags().StringP("match-only", "M", "", `Filters out results but only calculates up visible rows`)
	cmdObj.Flags().StringP("select", "", "", `Filters pods based on their spec field, comma seperated list of FIELD OP VALUE, where OP can be one of ==, = and != `)
	cmdObj.Flags().BoolP("show-namespace", "", false, `Show the namespace column`)
	cmdObj.Flags().BoolP("show-node", "", false, `Show the node name column`)
	cmdObj.Flags().BoolP("show-type", "T", false, `Show the container type column, where:
    I=init container, C=container, E=ephemerial container, P=Pod, D=Deployment, R=ReplicaSet, A=DaemonSet, S=StatefulSet, N=Node`)
	cmdObj.Flags().StringP("node-label", "", "", `Show the selected node label as a column`)
	cmdObj.Flags().StringP("pod-label", "", "", `Show the selected pod label as a column`)
	cmdObj.Flags().StringP("annotation", "", "", `Show the selected annotation as a column`)
	cmdObj.Flags().StringP("filename", "f", "", `read pod information from this yaml file instead`)
	cmdObj.Flags().StringP("columns", "", "", `list of column names to show in the table output, all other columns are hidden`)
	cmdObj.Flags().StringP("color", "", "", `Add some much needed colour to the table output. string can be one of: columns, custom, errors, mix and none (overrides env variable ICE_COLOUR)`)
}

func processCommonFlags(cmd *cobra.Command) (commonFlags, error) {
	var err error

	f := commonFlags{}

	if cmd.Flag("all-namespaces").Value.String() == "true" {
		f.allNamespaces = true
		f.showNamespaceName = true
	}

	if cmd.Flag("include-init") != nil {
		if cmd.Flag("include-init").Value.String() == "true" {
			f.showInitContainers = true
		}
	}

	if cmd.Flag("oddities") != nil {
		if cmd.Flag("oddities").Value.String() == "true" {
			f.showOddities = true
		}
	}

	if cmd.Flag("selector") != nil {
		if len(cmd.Flag("selector").Value.String()) > 0 {
			f.labels = cmd.Flag("selector").Value.String()
		}
	}

	if cmd.Flag("container") != nil {
		if len(cmd.Flag("container").Value.String()) > 0 {
			f.container = cmd.Flag("container").Value.String()
		}
	}

	if cmd.Flag("output") != nil {
		if len(cmd.Flag("output").Value.String()) > 0 {
			outAs := cmd.Flag("output").Value.String()
			// we use a switch to match -o flag so I can expand in future
			switch strings.ToLower(outAs) {
			case "csv":
				f.outputAs = "csv"
			case "list":
				f.outputAs = "list"
			case "json":
				f.outputAs = "json"
			case "yaml":
				f.outputAs = "yaml"

			default:
				return commonFlags{}, errors.New("unknown output format only csv, list, json and yaml are supported")
			}
		}
	}

	if cmd.Flag("size") != nil {
		if len(cmd.Flag("size").Value.String()) > 0 {
			f.byteSize = cmd.Flag("size").Value.String()
		}
	}

	if cmd.Flag("sort") != nil {
		// based on a whitelist approach sort just removes invalid chars,
		// we cant check header names as we dont know them at this point
		if len(cmd.Flag("sort").Value.String()) > 0 {
			rawSortString := cmd.Flag("sort").Value.String()
			f.sortList, err = splitAndFilterList(rawSortString, "ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-")
			if err != nil {
				return commonFlags{}, err
			}
		}
	}

	rawMatchString := ""
	if cmd.Flag("match") != nil {
		if len(cmd.Flag("match").Value.String()) > 0 {
			rawMatchString = cmd.Flag("match").Value.String()
		}
	}
	if cmd.Flag("match-only") != nil {
		if len(cmd.Flag("match-only").Value.String()) > 0 {
			rawMatchString = cmd.Flag("match-only").Value.String()
			f.calcMatchOnly = true
		}
	}
	if len(rawMatchString) > 0 {
		f.filterList, err = splitAndFilterMatchList(rawMatchString, "ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-.0123456789<>=*?", []string{"<=", ">=", "!=", "==", "=", "<", ">"})
		if err != nil {
			return commonFlags{}, err
		}
	}

	if cmd.Flag("tree") != nil {
		if cmd.Flag("tree").Value.String() == "true" {
			if len(f.sortList) != 0 {
				return commonFlags{}, errors.New("you may not use the tree and sort flags together")
			}
			f.showTreeView = true
		}
	}

	if cmd.Flag("node-tree") != nil {
		if cmd.Flag("node-tree").Value.String() == "true" {
			if len(f.sortList) != 0 {
				return commonFlags{}, errors.New("you may not use the node-tree and sort flags together")
			}
			f.showNodeTree = true
			f.showTreeView = true
		}
	}
	if cmd.Flag("select") != nil {
		if len(cmd.Flag("select").Value.String()) > 0 {
			rawFilterString := cmd.Flag("select").Value.String()
			f.matchSpecList, err = splitAndFilterMatchList(rawFilterString, "ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-0123456789<>=*?", []string{"!=", "==", "="})
			if err != nil {
				return commonFlags{}, err
			}
		}
	}

	if cmd.Flag("show-namespace").Value.String() == "true" {
		f.showNamespaceName = true
	}

	if cmd.Flag("show-node").Value.String() == "true" {
		f.showNodeName = true
	}

	if cmd.Flag("show-type").Value.String() == "true" {
		f.showContainerType = true
	}

	if cmd.Flag("node-label").Value.String() != "" {
		label := cmd.Flag("node-label").Value.String()
		f.labelNodeName = label
	}

	if cmd.Flag("pod-label").Value.String() != "" {
		label := cmd.Flag("pod-label").Value.String()
		f.labelPodName = label
	}

	if cmd.Flag("annotation").Value.String() != "" {
		annotation := cmd.Flag("annotation").Value.String()
		f.annotationPodName = annotation
	}

	if cmd.Flag("filename").Value.String() != "" {
		inputFilename := cmd.Flag("filename").Value.String()
		f.inputFilename = inputFilename
	}

	if cmd.Flag("columns").Value.String() != "" {
		f.showColumnByName = cmd.Flag("columns").Value.String()
	}

	// check and set coluring type to use, we also check for both spellings of colour
	colourOut := ""
	// check environment vars first
	colourOut = os.Getenv("ICE_COLOR")

	//then allow overiding with flags
	if cmd.Flag("color") != nil {
		if len(cmd.Flag("color").Value.String()) > 0 {
			colourOut = cmd.Flag("color").Value.String()
		}
	}

	if len(colourOut) > 0 {
		// we use a switch to match --colour flag so I can expand in future
		colourEnv := strings.ToLower(colourOut)
		colourSet := strings.Split(colourEnv, ";")

		switch strings.ToLower(colourSet[0]) {
		case "mix":
			f.outputAsColour = COLOUR_MIX
		case "columns":
			f.outputAsColour = COLOUR_COLUMNS
		case "errors":
			f.outputAsColour = COLOUR_ERRORS
		case "none":
			f.outputAsColour = COLOUR_NONE
		case "custom":
			// f.outputAsColour = COLOUR_CUSTOM
			f.useTheseColours, f.outputAsColour, err = getColourSetFromString(colourSet[1:])
			if err != nil {
				return commonFlags{}, err
			}

		default:
			return commonFlags{}, errors.New("unknown colour type only columns, custom, errors, mix and none are supported")
		}
	}

	return f, nil
}

func splitAndFilterList(rawSortString string, filterString string) ([]string, error) {
	// based on a whitelist approach sort just removes invalid chars,
	// we cant check header names as we dont know them at this point
	var sortList []string
	var rawCase string

	rawSortList := strings.Split(rawSortString, ",")
	for i := 0; i < len(rawSortList); i++ {
		safeStr := ""
		rawItem := strings.TrimSpace(rawSortList[i])
		if len(rawItem) <= 0 {
			continue
		}

		// current used chars in headers are A-Z ! and % nothing else is needed
		// so pointless using regex
		rawCase = strings.ToUpper(rawItem)
		for _, v := range strings.Split(rawCase, "") {
			if strings.Contains(filterString, v) {
				safeStr += v
			}
		}

		if len(safeStr) != len(rawItem) {
			return []string{}, errors.New("invalid characters in column name")
		}
		sortList = append(sortList, safeStr)
	}

	return sortList, nil
}

// splitAndFilterMatchList removes any chars not in filterList and splits the line based on values in []operator, returns a map[string]matchValue type.
//
//	the order of operatorList is important as the match is done on a first come first served basis
func splitAndFilterMatchList(rawSortString string, filterString string, operatorList []string) (map[string]matchValue, error) {
	// based on a whitelist approach sort just removes invalid chars,
	// we cant check header names as we dont know them at this point
	var rawCase string
	sortList := make(map[string]matchValue)

	rawSortList := strings.Split(rawSortString, ",")
	for i := 0; i < len(rawSortList); i++ {
		safeStr := ""
		rawItem := strings.TrimSpace(rawSortList[i])
		if len(rawItem) <= 0 {
			continue
		}

		for _, v := range strings.Split(rawItem, "") {
			rawCase = strings.ToUpper(v)
			if strings.Contains(filterString, rawCase) {
				safeStr += v
			}
		}

		if len(safeStr) != len(rawItem) {
			return map[string]matchValue{}, errors.New("invalid characters in suppiled string")
		}

		// find and split based on operatorList
		found := false
		fieldName := ""
		operator := ""
		value := ""

		for i := 0; i < len(operatorList); i++ {
			operator = operatorList[i]
			// check idx is 1 or more as we need at least a single charactor before the operator
			if idx := strings.Index(safeStr, operator); idx > 0 {
				fieldName = strings.ToUpper(strings.TrimSpace(safeStr[:idx]))
				value = strings.TrimSpace(safeStr[idx+len(operator):])
				found = true
				break
			}
		}

		if found {
			sortList[fieldName] = matchValue{
				operator: operator,
				value:    value,
			}
		}
	}

	return sortList, nil
}
