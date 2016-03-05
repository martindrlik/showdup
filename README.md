# showdup

Command displays identical files in specified directories.

## Usage

Run the `showdup` binary to display identical files in the current directory:

	$ showdup

Pass arguments to specify directories:

	$ showdup . dir1 dir2

Current implementation uses md5 checksum to find out if files are identical.
