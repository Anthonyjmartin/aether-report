package software_check

import (
	"github.com/shirou/gopsutil/docker"
)

func DockerDetails() (detailsStr []docker.CgroupDockerStat, err error) {
	detailsStr, err = docker.GetDockerStat()
	if err != nil {
		return
	}
	return
}
