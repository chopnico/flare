# flare

Flare is a CLI/TUI application for managing Cloudflare instances. The original intent was to build a TUI around Clouflare (because why not....), but I decided to extend it to make it more of an all encompassing command line utility. IaC is important, so I would reccomend using something like Terraform or Ansible over this tool, but eventually you'll be able to manage everything with flare. It could come in handy if you want to run ad-hoc commands such as "purge cache" or "enable development".

## Install

You can compile from source or you can grab the latest release.

[](https://github.com/chopnico/flare/releases)

## Support

Support will expand. (I'm a Linux guy, so priorities....)

* Linux AMD64

## How to Use

#### Show help
``` sh
flare --help
```

``` sh
NAME:
   flare - A Cloudflare CLI/CUI tool

USAGE:
   flare [global options] command [command options] [arguments...]

COMMANDS:
   init, i  Initialize configuration
   tui, t   Run the terminal UI
   zone, z  Interact with zones
   dns, d   Interact with DNS records
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output-format OUTPUT-FORMAT, -n OUTPUT-FORMAT  OUTPUT-FORMAT (default: "list")
   --help, -h                                       show help (default: false)
```

#### Show some zones

``` sh
flare zone list
```

``` sh
ID   : 00000000000000000000000000000000
Name : audia.cloud

ID   : 00000000000000000000000000000000
Name : audiacloud.me

ID   : 00000000000000000000000000000000
Name : gersh.win
```

#### Show some details about a zone

``` sh
flare zone detail --id 00000000000000000000000000000000
```

``` sh
ID                : 00000000000000000000000000000000
Name              : audia.cloud
DevMode           : 0
OriginalNS        : [ns2.hover.com ns1.hover.com]
OriginalRegistrar : <omit> 
OriginalDNSHost   : 
CreatedOn         : <omit>
ModifiedOn        : <omit> 
NameServers       : [cory.ns.cloudflare.com kami.ns.cloudflare.com]
Owner             : <omit>
Permissions       : <omit> 
PlanPending       : {}
Status            : active
Paused            : false
Type              : full
Host              : { }
VanityNS          : []
Betas             : []
DeactReason       : 
Meta              : {}
Account           : {}
VerificationKey   : 
```

#### Output details as JSON

``` sh
flare --output-format json zone detail --id 00000000000000000000000000000000
```

``` sh
you get the idea... all that json
```

