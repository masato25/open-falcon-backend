package falconPortal

import (
	"fmt"
	"time"

	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"

	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

const timeLayout = "2006-01-02 15:04:05"

//will deprecated
func CloseEvent(username string, colsed_note string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event_cases SET user_modified = ?, closed_at = ?, status = ?, closed_note = ? WHERE id = ?", userid, time.Now().Format(timeLayout), "SOLVED", colsed_note, id).Exec()
	return
}

func AddNote(username string, processNote string, eventcaseid string, processStatus string, bossId string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	eventcaseids := []string{}
	if strings.Contains(eventcaseid, ",") {
		eventcaseids = strings.Split(eventcaseid, ",")
	} else {
		eventcaseids = append(eventcaseids, eventcaseid)
	}
	for _, eventid := range eventcaseids {
		sqlbase := fmt.Sprintf("SET event_caseId = '%s' , user_id = %d", eventid, userid)
		if processNote != "" {
			sqlbase = fmt.Sprintf("%s , note = '%s'", sqlbase, processNote)
		}
		//for set boss case id
		if bossId != "" {
			sqlbase = fmt.Sprintf("%s, case_id = '%s'", sqlbase, bossId)
		}
		if processStatus != "" {
			sqlbase = fmt.Sprintf("%s, status = '%s'", sqlbase, processStatus)
		}
		var processNoteID int
		q.Begin()
		q.Raw(fmt.Sprintf("Insert INTO event_note %s, timestamp = ? ;", sqlbase), time.Now().Format(timeLayout)).Exec()
		err = q.Raw("SELECT LAST_INSERT_ID()").QueryRow(&processNoteID)
		if processNoteID != 0 && processStatus == "in progress" && err == nil {
			//update note status to event_cases
			_, err = q.Raw("Update event_cases SET process_note = ?, process_status = ? WHERE id = ?", processNoteID, processStatus, eventcaseid).Exec()
		} else if processNoteID != 0 && (processStatus == "resolved" || processStatus == "ignored") && err == nil {
			//update note status & add closed_at to event_cases
			_, err = q.Raw(fmt.Sprintf("Update event_cases SET process_note = %v, process_status = '%v', closed_at = FROM_UNIXTIME(%v) WHERE id = '%v'", processNoteID, processStatus, time.Now().Unix(), eventcaseid)).Exec()
		}
		q.Commit()
		if err != nil {
			return
		}
	}
	return
}

func IgnoredAlertTiemOutRecover() (rowaffect int64, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	resetBaseTime := time.Now().Unix() - g.Config().IgnoredResetTimout
	s, err := q.Raw("Update event_cases Set process_status = NULL, process_note = NULL, closed_at = NULL where closed_at < FROM_UNIXTIME(?) AND process_status = ?", resetBaseTime, "ignored").Exec()
	rowaffect, _ = s.RowsAffected()
	return
}
