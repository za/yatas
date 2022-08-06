package volumes

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetVolumes(s *session.Session) []*ec2.Volume {
	svc := ec2.New(s)
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(input)
	if err != nil {
		panic(err)
	}
	return result.Volumes
}

func checkIfEncryptionEnabled(s *session.Session, volumes []*ec2.Volume, c *[]types.Check) {
	logger.Info("Running AWS_VOL_001")
	var check types.Check
	check.Name = "EC2 Volumes Encryption"
	check.Id = "AWS_VOL_001"
	check.Description = "Check if EC2 encryption is enabled"
	check.Status = "OK"
	svc := ec2.New(s)
	for _, volume := range volumes {
		params := &ec2.DescribeVolumesInput{
			VolumeIds: []*string{volume.VolumeId},
		}
		resp, err := svc.DescribeVolumes(params)
		if err != nil {
			panic(err)
		}
		if *resp.Volumes[0].Encrypted == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfVolumesTypeGP3(s *session.Session, volumes []*ec2.Volume, c *[]types.Check) {
	logger.Info("Running AWS_VOL_004")
	var check types.Check
	check.Name = "EC2 Volumes Type"
	check.Id = "AWS_VOL_004"
	check.Description = "Check if all volumes are of type gp3"
	check.Status = "OK"
	for _, volume := range volumes {
		if *volume.VolumeType != "gp3" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunVolumesTest(s *session.Session) []types.Check {
	var checks []types.Check
	logger.Debug("Starting EC2 volumes tests")
	volumes := GetVolumes(s)
	snapshots := GetSnapshots(s)
	checkIfEncryptionEnabled(s, volumes, &checks)
	CheckIfAllVolumesHaveSnapshots(s, volumes, &checks)
	CheckIfAllSnapshotsEncrypted(s, snapshots, &checks)
	CheckIfVolumesTypeGP3(s, volumes, &checks)
	CheckIfSnapshotYoungerthan24h(s, snapshots, &checks)
	return checks
}
