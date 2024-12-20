package repositories

import (
	"fmt"
	"saymow/version-manager/app/repositories/filesystems"
	"testing"

	"gotest.tools/v3/fs"
)

func fixtureMakeBasicRepositoryFs(dir *fs.Dir) fs.PathOp {
	return fs.WithDir(
		filesystems.REPOSITORY_FOLDER_NAME,
		fs.WithDir(
			filesystems.SAVES_FOLDER_NAME,
			fs.WithFile(
				"9a35bd416196f27e40f4f9e4768496ef29c1922f0ab5e2651a218e4d4cb09688",
				fmt.Sprintf(`initial save

11/15 04:08:58PM '24 -0300
	
Please do not edit the lines below.
	
	
Files:
	
%s	(created)
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
%s	(created)
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`, dir.Join("1.txt"), dir.Join("2.txt")),
			),
			fs.WithFile(
				"3f674c71a3596db8f24fd31a85c503ae600898cc03810fcc171781d4f35531d2",
				fmt.Sprintf(`second save
9a35bd416196f27e40f4f9e4768496ef29c1922f0ab5e2651a218e4d4cb09688
11/15 04:09:54PM '24 -0300

Please do not edit the lines below.


Files:
	
%s	(created)
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
%s	(modified)
6f6367cbecfac86af4e749156e1b1046524eff9afbd8a29c964c3b46ebdf7fc2`, dir.Join("3.txt"), dir.Join("1.txt")),
			),
		),
		fs.WithDir(
			filesystems.OBJECTS_FOLDER_NAME,
			fs.WithFile("814f15a360c1a700342d1652e3bd8b9c954ee2ad9c974f6ec88eb92ff2d6b3b3", ""),
			fs.WithFile("6f6367cbecfac86af4e749156e1b1046524eff9afbd8a29c964c3b46ebdf7fc2", ""),
			fs.WithFile("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", ""),
		),
		fs.WithFile(filesystems.REFS_FILE_NAME, fmt.Sprintf("Refs:\r\n\r\n%s\r\n3f674c71a3596db8f24fd31a85c503ae600898cc03810fcc171781d4f35531d2\r\n", filesystems.INITIAL_REF_NAME)),
		fs.WithFile(filesystems.HEAD_FILE_NAME, filesystems.INITIAL_REF_NAME),
		fs.WithFile(filesystems.INDEX_FILE_NAME, fmt.Sprintf(`Tracked files:
	
%s	(created)
814f15a360c1a700342d1652e3bd8b9c954ee2ad9c974f6ec88eb92ff2d6b3b3
%s	(removed)`, dir.Join("4.txt"), dir.Join("2.txt"))),
	)
}

func fixtureGetBaseProject(t *testing.T) (*fs.Dir, *Repository) {
	dir := fs.NewDir(
		t,
		"project",
		fs.WithFile("1.txt", "1 content"),
		fs.WithFile("2.txt", "2 content"),
		fs.WithFile("3.txt", "3 content"),
		fs.WithDir(
			"a",
			fs.WithFile("4.txt", "4 content"),
			fs.WithFile("5.txt", "5 content"),
			fs.WithDir(
				"b",
				fs.WithFile("6.txt", "6 content"),
				fs.WithFile("7.txt", "7 content"),
			),
		),
		fs.WithDir(
			"c",
			fs.WithFile("8.txt", "8 content"),
			fs.WithFile("9.txt", "9 content"),
		),
	)

	return dir, CreateRepository(dir.Path())
}

func fixtureGetNewProject(t *testing.T) (*fs.Dir, *Repository) {
	dir := fs.NewDir(
		t,
		"project",
	)

	return dir, CreateRepository(dir.Path())
}

func fixtureGetCustomProject(t *testing.T, makeRepositoryDir func(dir *fs.Dir) fs.PathOp) (*fs.Dir, *Repository) {
	dir := fs.NewDir(
		t,
		"project",
		fs.WithFile("1.txt", "1 content"),
		fs.WithFile("2.txt", "2 content"),
		fs.WithFile("3.txt", "3 content"),
		fs.WithDir(
			"a",
			fs.WithFile("4.txt", "4 content"),
			fs.WithFile("5.txt", "5 content"),
			fs.WithDir(
				"b",
				fs.WithFile("6.txt", "6 content"),
				fs.WithFile("7.txt", "7 content"),
			),
		),
		fs.WithDir(
			"c",
			fs.WithFile("8.txt", "8 content"),
			fs.WithFile("9.txt", "9 content"),
		),
	)

	fs.Apply(t, dir, makeRepositoryDir(dir))

	return dir, GetRepository(dir.Path())
}
