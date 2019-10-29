#!/usr/local/python
# -*- coding: UTF-8 -*-

import argparse
import os
import shutil

# print("system: " + os.name)
# print(os.path.expanduser("~"))

if os.name == "nt":
    output = "bchaind.exe"
    install_path = os.path.expanduser("~") + "/bchain.io"
elif os.name == "posix":
    output = "bchaind"
    install_path = "/usr/local/bin/bchain.io"
else:
    print("unkown system")
    exit(0)

root_path = os.path.abspath(".")
out_path = root_path + "/bin/"
src_path = "bchain.io/bchaind"


def build(args):
    print("build ...")
    os.environ["GOPATH"] = root_path
    print(os.environ.get("GOPATH"))
    cmd = "go build -x -i -o " + out_path + output + " " + src_path
    print(cmd)
    if not os.path.exists(out_path):
        os.mkdir(out_path)
    os.system(cmd)


def install(args):
    print("install " + output + " to " + args.install_dir + " ...")

    if not os.path.exists(out_path + output):
        print("file: " + out_path + output + " is not exist.")
        exit(0)

    if args.install_dir != install_path:
        if not os.path.exists(args.install_dir):
            print("path: " + args.install_dir + " is not exist.")
            exit(0)

    if not os.path.exists(args.install_dir):
        os.mkdir(args.install_dir)

    shutil.copy(out_path + output, args.install_dir)


def main():
    parser = argparse.ArgumentParser()

    subparser = parser.add_subparsers(title='subcommands',
                                      # description='valid subcommands',
                                      help='additional help',
                                      dest="subcomand"
                                      )

    build_parser = subparser.add_parser("build", help="build bchain")
    build_parser.set_defaults(func=build)

    install_parser = subparser.add_parser("install", help="install bchain")
    install_parser.add_argument("-d", "--dir",
                                required=False,
                                type=str,
                                dest="install_dir",
                                metavar="PATH",
                                default=install_path,
                                help="install bchain to dir(default: %(default)s)")
    install_parser.set_defaults(func=install)

    # parser.add_argument("-b", "--build",
    #                     required=False,
    #                     action="store_true",
    #                     help="build bchain")
    #
    # parser.add_argument("-i", "--install",
    #                     required=False,
    #                     nargs="?",
    #                     # default=".",
    #                     dest="install_path",
    #                     metavar="PATH",
    #                     help="install bchain to path (default: .)")

    args = parser.parse_args()
    if args.subcomand:
        args.func(args)
    else:
        parser.print_help()

    return

if __name__ == "__main__":
    main()
