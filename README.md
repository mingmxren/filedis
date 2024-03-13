# filedis

filedis is a simple tool to do file distribution by file ext 
remain file tree structure.

## Usage

```shell
filedis --src <source> --dst ext1:dir1 --dst ext2:dir2 --dst ext3:dir3
```

## Example

```shell
# before
SOME_DIR
├── 001
│   ├── 001.JPG
│   └── 001.RAF
├── 002
│   ├── 002.JPG
│   └── 002.RAF
└── 003
    ├── 003.JPG
    └── 003.RAF

$ filedis --src SOME_DIR --dst raf:RAF_DIR

# after
SOME_DIR
├── 001
│   └── 001.JPG
├── 002
│   └── 002.JPG
└── 003
    └── 003.JPG
RAF_DIR
├── 001
│   └── 001.RAF
├── 002
│   └── 002.RAF
└── 003
    └── 003.RAF
```
