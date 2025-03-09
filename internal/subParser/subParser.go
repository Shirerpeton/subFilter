package subParser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type sub struct {
	original string
	dialogs []dialog
}

type dialog struct {
	start time.Duration
	end time.Duration
	content string
}

func getDurationFromTimestamp(timestamp string) (time.Duration, error) {
	parts := strings.Split(timestamp, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in hours: %v", err)
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in minutes: %v", err)
		return 0, err
	}
	var sep string
	if strings.Contains(parts[2], ".") {
		sep = "."
	} else if strings.Contains(parts[2], ",") {
		sep = ","
	}
	parts = strings.Split(parts[2], sep)
	seconds, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in seconds: %v", err)
		return 0, err
	}
	if len(parts[1]) == 2 {
		parts[1] += "0"
	}
	miliseconds, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in miliseconds: %v", err)
		return 0, err
	}
	duration := hours*int(time.Hour) + minutes*int(time.Minute) + seconds*int(time.Second) + miliseconds*int(time.Millisecond);
	return time.Duration(duration), nil
}

func getSrtTimestampFromDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	miliseconds := int(d.Milliseconds()) % 1000
	timestamp := fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, miliseconds)
	return timestamp
}

func GetSubFromSrt(content string) (sub, error) {
	s := sub{
		original: content,
	}
	ds := make([]dialog, 0)
	d := &dialog{}
	for line := range strings.SplitSeq(content, "\n") {
		if strings.Contains(line, "-->") {
			parts := strings.Split(line, " ")
			start, err := getDurationFromTimestamp(parts[0])
			if err != nil {
				return s, err
			}
			end, err := getDurationFromTimestamp(parts[2])
			if err != nil {
				return s, err
			}
			d = &dialog{
				start: start,
				end: end,
			}
		} else {
			if d == nil {
				continue
			}
			if line == "" {
				ds = append(ds, *d)
				d = nil
				continue
			}
			if d.content == "" {
				d.content = line
			} else {
				d.content += "\n" + line
			}
		}
	}
	s.dialogs = ds
	return s, nil
}

func (s *sub) Filter(str string) {
	ds := make([]dialog, 0, len(s.dialogs))
	for _, d := range s.dialogs {
		if !strings.Contains(d.content, str) {
			ds = append(ds, d)
		}
	}
	s.dialogs = ds
}

func (s sub) GetSrtFromSub() string {
	b := strings.Builder{}
	for i, d := range s.dialogs {
		b.WriteString(fmt.Sprintf("%d\n", i + 1))
		b.WriteString(fmt.Sprintf("%s --> %s\n", getSrtTimestampFromDuration(d.start), getSrtTimestampFromDuration(d.end)))
		b.WriteString(d.content)
		if len(s.dialogs) - 1 != i {
			b.WriteString("\n\n")
		}
	}
	return b.String()
}
