package models

import (
	"encoding/csv"
	"github.com/cSploit/daemon/models/internal"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

var aircrackTimeLayout = "2006-01-02 15:04:05"

type DiscoveryJob struct {
	internal.Base

	Dir string `json:"-"`

	Job   Job
	JobId uint
}

func (d *DiscoveryJob) parseOne(file string) error {
	// Dirty hack to have a clean dump
	//TODO: from stdout ?
	dump, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	dump_str := string(dump)
	// Replace endline with just an  \n
	dump_str = strings.Replace(dump_str, ", \r\n", ", \n", -1)
	dump_str = strings.Replace(dump_str, ",\r\n", ",\n", -1)
	dump_split := strings.SplitN(dump_str, "\r\n", 4)

	// Extract the two parts of the csv
	dump_aps := dump_split[2]
	dump_clients := dump_split[3]
	dump_clients = strings.SplitN(dump_clients, "\r\n", 2)[1]

	// End of dirty hack, fill the structs
	reader_aps := csv.NewReader(strings.NewReader(dump_aps))
	reader_clients := csv.NewReader(strings.NewReader(dump_clients))

	// Start with the aps
	for {
		record, csv_err := reader_aps.Read()
		if csv_err == io.EOF {
			break
		}
		if csv_err != nil {
			return err
		}

		// Okay, fill an AP struct then append to the dump
		ap, e := FindApByBssid(record[0])

		if e == gorm.ErrRecordNotFound {
			ap = &AP{}
		} else if e != nil {
			log.Error(e)
			continue
		}

		// TODO: clean that
		// FIXME: I am too lazy to check the errors

		ap.Bssid = record[0]
		ap.First, _ = time.Parse(aircrackTimeLayout, record[1])
		ap.Last, _ = time.Parse(aircrackTimeLayout, record[2])
		ap.Channel, _ = strconv.Atoi(record[3])
		ap.Speed, _ = strconv.Atoi(record[4])
		ap.Privacy = record[5]
		ap.Cipher = record[6]
		ap.Auth = record[7]
		ap.Power, _ = strconv.Atoi(record[8])
		ap.Beacons, _ = strconv.Atoi(record[9])
		ap.IVs, _ = strconv.Atoi(record[10])
		ap.Lan = strings.Replace(record[11], " ", "", -1)
		ap.IdLen, _ = strconv.Atoi(record[12])
		ap.Essid = record[13]
		ap.Key = record[14]

		if err := internal.Db.Save(ap); err != nil {
			log.Error(err)
		} else if err := internal.Db.Model(&(d.Job)).Association("job_aps").Append(ap).Error; err != nil {
			log.Error(err)
		}
	}

	// Continue with the clients
	for {
		record, csv_err := reader_clients.Read()
		if csv_err == io.EOF {
			break
		}
		if csv_err != nil {
			return err
		}

		// Okay, fill a Client struct then append to the dump
		client, e := FindClientByMac(record[0])

		if e == gorm.ErrRecordNotFound {
			client = &Client{}
		} else if e != nil {
			log.Error(e)
			continue
		}

		// TODO: clean that
		// FIXME: too lazy to fix the errors
		client.Station = record[0]
		client.First, _ = time.Parse(aircrackTimeLayout, record[1])
		client.Last, _ = time.Parse(aircrackTimeLayout, record[2])
		client.Power, _ = strconv.Atoi(record[3])
		client.Packets, _ = strconv.Atoi(record[4])
		client.Bssid = record[5]
		client.Probed = record[6]

		if err := internal.Db.Save(client); err != nil {
			log.Error(err)
		} else if err := internal.Db.Model(&(d.Job)).Association("job_clients").Append(client).Error; err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (d *DiscoveryJob) Parse() error {
	files, e := ioutil.ReadDir(d.Dir)

	if e != nil {
		return e
	}

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".csv") {
			continue
		}

		if err := d.parseOne(path.Join(d.Dir, fi.Name())); err != nil {
			return err
		}
	}

	return nil
}
