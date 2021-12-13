# Tool Documentation (create-gp4)

## Summary
`create-gp4` takes a content ID and list of files and produces an output GP4 project file. This tool is based on John Tornblom's `create-gp4` prototype python script.

### Build
The tool is written in Golang. Commands should be the same between Windows and Linux assuming Golang is installed on the target system.

Building is straightforward, simply go to the root directory of the tool (`/src/tools/create-gp4/`) and run the `go build` command.

There is a shell and batch build script to compile `create-gp4` on all platforms (Windows, Linux, and macOS).

### Usage
The `create-gp4` tool requires two arguments. These arguments are `-content-id` and `-files` respectively. The `-out` path for the gp4 file is optional, and defaults to `homebrew.gp4`.

```
Usage of create-gp4:
  -out string
        output gp4 to write to
  -content-id string
        content ID of the package
  -files string
        list of files to pack into the package
```

### Example Usage

```
./create-gp4 -out IV0000-BREW00083_00-HELLOWORLD000000.gp4 --content-id=IV0000-BREW00083_00-HELLOWORLD000000 --files eboot.bin sce_sys/about/right.sprx sce_sys/param.sfo sce_sys/icon0.png sce_module/libSceFios2.prx sce_module/libc.prx
```

## Source code overview
All the source code resides in `main.go` since the program is relatively trivial.
