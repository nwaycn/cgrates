/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or56
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package engine

import (
	"reflect"
	"testing"

	"github.com/cgrates/cgrates/utils"
)

func TestModelHelperCsvLoad(t *testing.T) {
	l, err := csvLoad(TpDestination{}, []string{"TEST_DEST", "+492"})
	tpd, ok := l.(TpDestination)
	if err != nil || !ok || tpd.Tag != "TEST_DEST" || tpd.Prefix != "+492" {
		t.Errorf("model load failed: %+v", tpd)
	}
}

func TestModelHelperCsvLoadInt(t *testing.T) {
	l, err := csvLoad(TpCdrstat{}, []string{"CDRST1", "5", "60m", "10s", "ASR", "2014-07-29T15:00:00Z;2014-07-29T16:00:00Z", "*voice", "87.139.12.167", "FS_JSON", "*rated", "*out", "cgrates.org", "call", "dan", "dan", "49", "3m;7m", "5m;10m", "suppl1", "NORMAL_CLEARING", "default", "rif", "rif", "0;2", "STANDARD_TRIGGERS"})
	tpd, ok := l.(TpCdrstat)
	if err != nil || !ok || tpd.QueueLength != 5 {
		t.Errorf("model load failed: %+v", tpd)
	}
}

func TestModelHelperCsvDump(t *testing.T) {
	tpd := TpDestination{
		Tag:    "TEST_DEST",
		Prefix: "+492"}
	csv, err := csvDump(tpd)
	if err != nil || csv[0] != "TEST_DEST" || csv[1] != "+492" {
		t.Errorf("model load failed: %+v", tpd)
	}
}

func TestTPDestinationAsExportSlice(t *testing.T) {
	tpDst := &utils.TPDestination{
		TPid:     "TEST_TPID",
		ID:       "TEST_DEST",
		Prefixes: []string{"49", "49176", "49151"},
	}
	expectedSlc := [][]string{
		[]string{"TEST_DEST", "49"},
		[]string{"TEST_DEST", "49176"},
		[]string{"TEST_DEST", "49151"},
	}
	mdst := APItoModelDestination(tpDst)
	var slc [][]string
	for _, md := range mdst {
		lc, err := csvDump(md)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTpDestinationsAsTPDestinations(t *testing.T) {
	tpd1 := TpDestination{Tpid: "TEST_TPID", Tag: "TEST_DEST", Prefix: "+491"}
	tpd2 := TpDestination{Tpid: "TEST_TPID", Tag: "TEST_DEST", Prefix: "+492"}
	tpd3 := TpDestination{Tpid: "TEST_TPID", Tag: "TEST_DEST", Prefix: "+493"}
	eTPDestinations := []*utils.TPDestination{&utils.TPDestination{TPid: "TEST_TPID", ID: "TEST_DEST",
		Prefixes: []string{"+491", "+492", "+493"}}}
	if tpDst := TpDestinations([]TpDestination{tpd1, tpd2, tpd3}).AsTPDestinations(); !reflect.DeepEqual(eTPDestinations, tpDst) {
		t.Errorf("Expecting: %+v, received: %+v", eTPDestinations, tpDst)
	}

}

func TestTPRateAsExportSlice(t *testing.T) {
	tpRate := &utils.TPRate{
		TPid: "TEST_TPID",
		ID:   "TEST_RATEID",
		RateSlots: []*utils.RateSlot{
			&utils.RateSlot{
				ConnectFee:         0.100,
				Rate:               0.200,
				RateUnit:           "60",
				RateIncrement:      "60",
				GroupIntervalStart: "0"},
			&utils.RateSlot{
				ConnectFee:         0.0,
				Rate:               0.1,
				RateUnit:           "1",
				RateIncrement:      "60",
				GroupIntervalStart: "60"},
		},
	}
	expectedSlc := [][]string{
		[]string{"TEST_RATEID", "0.1", "0.2", "60", "60", "0"},
		[]string{"TEST_RATEID", "0", "0.1", "1", "60", "60"},
	}

	ms := APItoModelRate(tpRate)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc[0], slc[0])
	}
}

func TestTPDestinationRateAsExportSlice(t *testing.T) {
	tpDstRate := &utils.TPDestinationRate{
		TPid: "TEST_TPID",
		ID:   "TEST_DSTRATE",
		DestinationRates: []*utils.DestinationRate{
			&utils.DestinationRate{
				DestinationId:    "TEST_DEST1",
				RateId:           "TEST_RATE1",
				RoundingMethod:   "*up",
				RoundingDecimals: 4},
			&utils.DestinationRate{
				DestinationId:    "TEST_DEST2",
				RateId:           "TEST_RATE2",
				RoundingMethod:   "*up",
				RoundingDecimals: 4},
		},
	}
	expectedSlc := [][]string{
		[]string{"TEST_DSTRATE", "TEST_DEST1", "TEST_RATE1", "*up", "4", "0", ""},
		[]string{"TEST_DSTRATE", "TEST_DEST2", "TEST_RATE2", "*up", "4", "0", ""},
	}
	ms := APItoModelDestinationRate(tpDstRate)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}

	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}

}

func TestApierTPTimingAsExportSlice(t *testing.T) {
	tpTiming := &utils.ApierTPTiming{
		TPid:      "TEST_TPID",
		ID:        "TEST_TIMING",
		Years:     "*any",
		Months:    "*any",
		MonthDays: "*any",
		WeekDays:  "1;2;4",
		Time:      "00:00:01"}
	expectedSlc := [][]string{
		[]string{"TEST_TIMING", "*any", "*any", "*any", "1;2;4", "00:00:01"},
	}
	ms := APItoModelTiming(tpTiming)
	var slc [][]string

	lc, err := csvDump(ms)
	if err != nil {
		t.Error("Error dumping to csv: ", err)
	}
	slc = append(slc, lc)

	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

/*
func TestAPItoModelStats(t *testing.T) {
	tpS := &utils.TPStats{
		TPid: "TPS1",
		ID:   "Stat1",
		Filters: []*utils.TPRequestFilter{
			&utils.TPRequestFilter{
				Type:      "*string",
				FieldName: "Account",
				Values:    []string{"1002"},
			},
		},
		ActivationInterval: &utils.TPActivationInterval{
			ActivationTime: "2014-07-29T15:00:00Z",
			ExpiryTime:     "",
		},
		TTL:        "1",
		Metrics:    []string{"MetricValue"},
		Blocker:    true,
		Stored:     true,
		Weight:     20,
		Thresholds: nil,
	}
	expectedSlc := [][]string{
		[]string{,"TPS1", "*Stat1", "*string", "*Account", "1002", "2014-07-29T15:00:00Z","","1","MetricValue",},
	}
	expectedtpS := APItoModelStats(tpS)
	var slc [][]string
	lc, err := csvDump(expectedtpS)
	if err != nil {
		t.Error("Error dumping to csv: ", err)
	}
	slc = append(slc, lc)

	if !reflect.DeepEqual(expectedtpS, tpS) {
		t.Errorf("Expecting: %+v, received: %+v", expectedtpS, slc)
	}
}

*/

func TestTPRatingPlanAsExportSlice(t *testing.T) {
	tpRpln := &utils.TPRatingPlan{
		TPid: "TEST_TPID",
		ID:   "TEST_RPLAN",
		RatingPlanBindings: []*utils.TPRatingPlanBinding{
			&utils.TPRatingPlanBinding{
				DestinationRatesId: "TEST_DSTRATE1",
				TimingId:           "TEST_TIMING1",
				Weight:             10.0},
			&utils.TPRatingPlanBinding{
				DestinationRatesId: "TEST_DSTRATE2",
				TimingId:           "TEST_TIMING2",
				Weight:             20.0},
		}}
	expectedSlc := [][]string{
		[]string{"TEST_RPLAN", "TEST_DSTRATE1", "TEST_TIMING1", "10"},
		[]string{"TEST_RPLAN", "TEST_DSTRATE2", "TEST_TIMING2", "20"},
	}

	ms := APItoModelRatingPlan(tpRpln)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTPRatingProfileAsExportSlice(t *testing.T) {
	tpRpf := &utils.TPRatingProfile{
		TPid:      "TEST_TPID",
		LoadId:    "TEST_LOADID",
		Direction: utils.OUT,
		Tenant:    "cgrates.org",
		Category:  "call",
		Subject:   "*any",
		RatingPlanActivations: []*utils.TPRatingActivation{
			&utils.TPRatingActivation{
				ActivationTime:   "2014-01-14T00:00:00Z",
				RatingPlanId:     "TEST_RPLAN1",
				FallbackSubjects: "subj1;subj2"},
			&utils.TPRatingActivation{
				ActivationTime:   "2014-01-15T00:00:00Z",
				RatingPlanId:     "TEST_RPLAN2",
				FallbackSubjects: "subj1;subj2"},
		},
	}
	expectedSlc := [][]string{
		[]string{utils.OUT, "cgrates.org", "call", "*any", "2014-01-14T00:00:00Z", "TEST_RPLAN1", "subj1;subj2", ""},
		[]string{utils.OUT, "cgrates.org", "call", "*any", "2014-01-15T00:00:00Z", "TEST_RPLAN2", "subj1;subj2", ""},
	}

	ms := APItoModelRatingProfile(tpRpf)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}

	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTPActionsAsExportSlice(t *testing.T) {
	tpActs := &utils.TPActions{
		TPid: "TEST_TPID",
		ID:   "TEST_ACTIONS",
		Actions: []*utils.TPAction{
			&utils.TPAction{
				Identifier:      "*topup_reset",
				BalanceType:     "*monetary",
				Directions:      utils.OUT,
				Units:           "5.0",
				ExpiryTime:      "*never",
				DestinationIds:  "*any",
				RatingSubject:   "special1",
				Categories:      "call",
				SharedGroups:    "GROUP1",
				BalanceWeight:   "10.0",
				ExtraParameters: "",
				Weight:          10.0},
			&utils.TPAction{
				Identifier:      "*http_post",
				BalanceType:     "",
				Directions:      "",
				Units:           "0.0",
				ExpiryTime:      "",
				DestinationIds:  "",
				RatingSubject:   "",
				Categories:      "",
				SharedGroups:    "",
				BalanceWeight:   "0.0",
				ExtraParameters: "http://localhost/&param1=value1",
				Weight:          20.0},
		},
	}
	expectedSlc := [][]string{
		[]string{"TEST_ACTIONS", "*topup_reset", "", "", "", "*monetary", utils.OUT, "call", "*any", "special1", "GROUP1", "*never", "", "5.0", "10.0", "", "", "10"},
		[]string{"TEST_ACTIONS", "*http_post", "http://localhost/&param1=value1", "", "", "", "", "", "", "", "", "", "", "0.0", "0.0", "", "", "20"},
	}

	ms := APItoModelAction(tpActs)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}

	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: \n%+v, received: \n%+v", expectedSlc, slc)
	}
}

// SHARED_A,*any,*highest,
func TestTPSharedGroupsAsExportSlice(t *testing.T) {
	tpSGs := &utils.TPSharedGroups{
		TPid: "TEST_TPID",
		ID:   "SHARED_GROUP_TEST",
		SharedGroups: []*utils.TPSharedGroup{
			&utils.TPSharedGroup{
				Account:       "*any",
				Strategy:      "*highest",
				RatingSubject: "special1"},
			&utils.TPSharedGroup{
				Account:       "second",
				Strategy:      "*highest",
				RatingSubject: "special2"},
		},
	}
	expectedSlc := [][]string{
		[]string{"SHARED_GROUP_TEST", "*any", "*highest", "special1"},
		[]string{"SHARED_GROUP_TEST", "second", "*highest", "special2"},
	}

	ms := APItoModelSharedGroup(tpSGs)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

//*in,cgrates.org,*any,EU_LANDLINE,LCR_STANDARD,*static,ivo;dan;rif,2012-01-01T00:00:00Z,10
func TestTPLcrRulesAsExportSlice(t *testing.T) {
	lcr := &utils.TPLcrRules{
		TPid:      "TEST_TPID",
		Direction: "*in",
		Tenant:    "cgrates.org",
		Category:  "LCR_STANDARD",
		Account:   "*any",
		Subject:   "*any",
		Rules: []*utils.TPLcrRule{
			&utils.TPLcrRule{
				DestinationId:  "EU_LANDLINE",
				Strategy:       "*static",
				StrategyParams: "ivo;dan;rif",
				ActivationTime: "2012-01-01T00:00:00Z",
				Weight:         20.0},
			//*in,cgrates.org,*any,*any,LCR_STANDARD,*lowest_cost,,2012-01-01T00:00:00Z,20
			&utils.TPLcrRule{
				DestinationId:  "*any",
				Strategy:       "*lowest_cost",
				StrategyParams: "",
				ActivationTime: "2012-01-01T00:00:00Z",
				Weight:         10.0},
		},
	}
	expectedSlc := [][]string{
		[]string{"*in", "cgrates.org", "LCR_STANDARD", "*any", "*any", "EU_LANDLINE", "", "*static", "ivo;dan;rif", "2012-01-01T00:00:00Z", "20"},
		[]string{"*in", "cgrates.org", "LCR_STANDARD", "*any", "*any", "*any", "", "*lowest_cost", "", "2012-01-01T00:00:00Z", "10"},
	}
	ms := APItoModelLcrRule(lcr)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

//CDRST1,5,60m,ASR,2014-07-29T15:00:00Z;2014-07-29T16:00:00Z,*voice,87.139.12.167,FS_JSON,rated,*out,cgrates.org,call,dan,dan,49,5m;10m,default,rif,rif,0;2,STANDARD_TRIGGERS
func TestTPCdrStatsAsExportSlice(t *testing.T) {
	cdrStats := &utils.TPCdrStats{
		TPid: "TEST_TPID",
		ID:   "CDRST1",
		CdrStats: []*utils.TPCdrStat{
			&utils.TPCdrStat{
				QueueLength:      "5",
				TimeWindow:       "60m",
				SaveInterval:     "10s",
				Metrics:          "ASR;ACD",
				SetupInterval:    "2014-07-29T15:00:00Z;2014-07-29T16:00:00Z",
				TORs:             "*voice",
				CdrHosts:         "87.139.12.167",
				CdrSources:       "FS_JSON",
				ReqTypes:         utils.META_RATED,
				Directions:       "*out",
				Tenants:          "cgrates.org",
				Categories:       "call",
				Accounts:         "dan",
				Subjects:         "dan",
				DestinationIds:   "49",
				PddInterval:      "3m;7m",
				UsageInterval:    "5m;10m",
				Suppliers:        "supplier1",
				DisconnectCauses: "NORMAL_CLEARNING",
				MediationRunIds:  "default",
				RatedAccounts:    "rif",
				RatedSubjects:    "rif",
				CostInterval:     "0;2",
				ActionTriggers:   "STANDARD_TRIGGERS"},
			&utils.TPCdrStat{
				QueueLength:      "5",
				TimeWindow:       "60m",
				SaveInterval:     "9s",
				Metrics:          "ASR",
				SetupInterval:    "2014-07-29T15:00:00Z;2014-07-29T16:00:00Z",
				TORs:             "*voice",
				CdrHosts:         "87.139.12.167",
				CdrSources:       "FS_JSON",
				ReqTypes:         utils.META_RATED,
				Directions:       "*out",
				Tenants:          "cgrates.org",
				Categories:       "call",
				Accounts:         "dan",
				Subjects:         "dan",
				DestinationIds:   "49",
				PddInterval:      "3m;7m",
				UsageInterval:    "5m;10m",
				Suppliers:        "supplier1",
				DisconnectCauses: "NORMAL_CLEARNING",
				MediationRunIds:  "default",
				RatedAccounts:    "dan",
				RatedSubjects:    "dan",
				CostInterval:     "0;2",
				ActionTriggers:   "STANDARD_TRIGGERS"},
		},
	}
	expectedSlc := [][]string{
		[]string{"CDRST1", "5", "60m", "10s", "ASR;ACD", "2014-07-29T15:00:00Z;2014-07-29T16:00:00Z", "*voice", "87.139.12.167", "FS_JSON", utils.META_RATED, "*out", "cgrates.org", "call",
			"dan", "dan", "49", "3m;7m", "5m;10m", "supplier1", "NORMAL_CLEARNING", "default", "rif", "rif", "0;2", "STANDARD_TRIGGERS"},
		[]string{"CDRST1", "5", "60m", "9s", "ASR", "2014-07-29T15:00:00Z;2014-07-29T16:00:00Z", "*voice", "87.139.12.167", "FS_JSON", utils.META_RATED, "*out", "cgrates.org", "call",
			"dan", "dan", "49", "3m;7m", "5m;10m", "supplier1", "NORMAL_CLEARNING", "default", "dan", "dan", "0;2", "STANDARD_TRIGGERS"},
	}
	ms := APItoModelCdrStat(cdrStats)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

//#Direction,Tenant,Category,Account,Subject,RunId,RunFilter,ReqTypeField,DirectionField,TenantField,CategoryField,AccountField,SubjectField,DestinationField,SetupTimeField,AnswerTimeField,UsageField
//*out,cgrates.org,call,1001,1001,derived_run1,,^rated,*default,*default,*default,*default,^1002,*default,*default,*default,*default
func TestTPDerivedChargersAsExportSlice(t *testing.T) {
	dcs := &utils.TPDerivedChargers{
		TPid:      "TEST_TPID",
		LoadId:    "TEST_LOADID",
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "1001",
		Subject:   "1001",
		DerivedChargers: []*utils.TPDerivedCharger{
			&utils.TPDerivedCharger{
				RunId:                "derived_run1",
				RunFilters:           "",
				ReqTypeField:         "^rated",
				DirectionField:       utils.META_DEFAULT,
				TenantField:          utils.META_DEFAULT,
				CategoryField:        utils.META_DEFAULT,
				AccountField:         utils.META_DEFAULT,
				SubjectField:         "^1002",
				DestinationField:     utils.META_DEFAULT,
				SetupTimeField:       utils.META_DEFAULT,
				PddField:             utils.META_DEFAULT,
				AnswerTimeField:      utils.META_DEFAULT,
				UsageField:           utils.META_DEFAULT,
				SupplierField:        utils.META_DEFAULT,
				DisconnectCauseField: utils.META_DEFAULT,
				CostField:            utils.META_DEFAULT,
				RatedField:           utils.META_DEFAULT,
			},
			&utils.TPDerivedCharger{
				RunId:                "derived_run2",
				RunFilters:           "",
				ReqTypeField:         "^rated",
				DirectionField:       utils.META_DEFAULT,
				TenantField:          utils.META_DEFAULT,
				CategoryField:        utils.META_DEFAULT,
				AccountField:         "^1002",
				SubjectField:         utils.META_DEFAULT,
				DestinationField:     utils.META_DEFAULT,
				SetupTimeField:       utils.META_DEFAULT,
				PddField:             utils.META_DEFAULT,
				AnswerTimeField:      utils.META_DEFAULT,
				UsageField:           utils.META_DEFAULT,
				SupplierField:        utils.META_DEFAULT,
				DisconnectCauseField: utils.META_DEFAULT,
				RatedField:           utils.META_DEFAULT,
				CostField:            utils.META_DEFAULT,
			},
		},
	}
	expectedSlc := [][]string{
		[]string{"*out", "cgrates.org", "call", "1001", "1001", "",
			"derived_run1", "", "^rated", utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, "^1002", utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT},
		[]string{"*out", "cgrates.org", "call", "1001", "1001", "",
			"derived_run2", "", "^rated", utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, "^1002", utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT, utils.META_DEFAULT},
	}
	ms := APItoModelDerivedCharger(dcs)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTPActionTriggersAsExportSlice(t *testing.T) {
	ap := &utils.TPActionPlan{
		TPid: "TEST_TPID",
		ID:   "PACKAGE_10",
		ActionPlan: []*utils.TPActionTiming{
			&utils.TPActionTiming{
				ActionsId: "TOPUP_RST_10",
				TimingId:  "ASAP",
				Weight:    10.0},
			&utils.TPActionTiming{
				ActionsId: "TOPUP_RST_5",
				TimingId:  "ASAP",
				Weight:    20.0},
		},
	}
	expectedSlc := [][]string{
		[]string{"PACKAGE_10", "TOPUP_RST_10", "ASAP", "10"},
		[]string{"PACKAGE_10", "TOPUP_RST_5", "ASAP", "20"},
	}
	ms := APItoModelActionPlan(ap)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTPActionPlanAsExportSlice(t *testing.T) {
	at := &utils.TPActionTriggers{
		TPid: "TEST_TPID",
		ID:   "STANDARD_TRIGGERS",
		ActionTriggers: []*utils.TPActionTrigger{
			&utils.TPActionTrigger{
				Id:                    "STANDARD_TRIGGERS",
				UniqueID:              "1",
				ThresholdType:         "*min_balance",
				ThresholdValue:        2.0,
				Recurrent:             false,
				MinSleep:              "0",
				BalanceId:             "b1",
				BalanceType:           "*monetary",
				BalanceDirections:     "*out",
				BalanceDestinationIds: "",
				BalanceWeight:         "0.0",
				BalanceExpirationDate: "*never",
				BalanceTimingTags:     "T1",
				BalanceRatingSubject:  "special1",
				BalanceCategories:     "call",
				BalanceSharedGroups:   "SHARED_1",
				BalanceBlocker:        "false",
				BalanceDisabled:       "false",
				MinQueuedItems:        0,
				ActionsId:             "LOG_WARNING",
				Weight:                10},
			&utils.TPActionTrigger{
				Id:                    "STANDARD_TRIGGERS",
				UniqueID:              "2",
				ThresholdType:         "*max_event_counter",
				ThresholdValue:        5.0,
				Recurrent:             false,
				MinSleep:              "0",
				BalanceId:             "b2",
				BalanceType:           "*monetary",
				BalanceDirections:     "*out",
				BalanceDestinationIds: "FS_USERS",
				BalanceWeight:         "0.0",
				BalanceExpirationDate: "*never",
				BalanceTimingTags:     "T1",
				BalanceRatingSubject:  "special1",
				BalanceCategories:     "call",
				BalanceSharedGroups:   "SHARED_1",
				BalanceBlocker:        "false",
				BalanceDisabled:       "false",
				MinQueuedItems:        0,
				ActionsId:             "LOG_WARNING",
				Weight:                10},
		},
	}
	expectedSlc := [][]string{
		[]string{"STANDARD_TRIGGERS", "1", "*min_balance", "2", "false", "0", "", "", "b1", "*monetary", "*out", "call", "", "special1", "SHARED_1", "*never", "T1", "0.0", "false", "false", "0", "LOG_WARNING", "10"},
		[]string{"STANDARD_TRIGGERS", "2", "*max_event_counter", "5", "false", "0", "", "", "b2", "*monetary", "*out", "call", "FS_USERS", "special1", "SHARED_1", "*never", "T1", "0.0", "false", "false", "0", "LOG_WARNING", "10"},
	}
	ms := APItoModelActionTrigger(at)
	var slc [][]string
	for _, m := range ms {
		lc, err := csvDump(m)
		if err != nil {
			t.Error("Error dumping to csv: ", err)
		}
		slc = append(slc, lc)
	}
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTPAccountActionsAsExportSlice(t *testing.T) {
	aa := &utils.TPAccountActions{
		TPid:             "TEST_TPID",
		LoadId:           "TEST_LOADID",
		Tenant:           "cgrates.org",
		Account:          "1001",
		ActionPlanId:     "PACKAGE_10_SHARED_A_5",
		ActionTriggersId: "STANDARD_TRIGGERS",
	}
	expectedSlc := [][]string{
		[]string{"cgrates.org", "1001", "PACKAGE_10_SHARED_A_5", "STANDARD_TRIGGERS", "false", "false"},
	}
	ms := APItoModelAccountAction(aa)
	var slc [][]string
	lc, err := csvDump(*ms)
	if err != nil {
		t.Error("Error dumping to csv: ", err)
	}
	slc = append(slc, lc)
	if !reflect.DeepEqual(expectedSlc, slc) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSlc, slc)
	}
}

func TestTpResourcesAsTpResources(t *testing.T) {
	tps := []*TpResource{
		&TpResource{
			Tpid:               "TEST_TPID",
			Tag:                "ResGroup1",
			FilterType:         MetaStringPrefix,
			FilterFieldName:    "Destination",
			FilterFieldValues:  "+49151;+49161",
			ActivationInterval: "2014-07-29T15:00:00Z",
			Stored:             false,
			Blocker:            false,
			Weight:             10.0,
			Limit:              "45",
			Thresholds:         "WARN_RES1;WARN_RES2"},
		&TpResource{
			Tpid:              "TEST_TPID",
			Tag:               "ResGroup1",
			FilterType:        MetaStringPrefix,
			FilterFieldName:   "Category",
			FilterFieldValues: "call;inbound_call",
			Thresholds:        "WARN3"},
		&TpResource{
			Tpid:               "TEST_TPID",
			Tag:                "ResGroup2",
			FilterType:         MetaStringPrefix,
			FilterFieldName:    "Destination",
			FilterFieldValues:  "+40",
			ActivationInterval: "2014-07-29T15:00:00Z",
			Stored:             false,
			Blocker:            false,
			Weight:             10.0,
			Limit:              "20"},
	}
	eTPs := []*utils.TPResource{
		&utils.TPResource{
			TPid: tps[0].Tpid,
			ID:   tps[0].Tag,
			Filters: []*utils.TPRequestFilter{
				&utils.TPRequestFilter{
					Type:      tps[0].FilterType,
					FieldName: tps[0].FilterFieldName,
					Values:    []string{"+49151", "+49161"},
				},
				&utils.TPRequestFilter{
					Type:      tps[1].FilterType,
					FieldName: tps[1].FilterFieldName,
					Values:    []string{"call", "inbound_call"},
				},
			},
			ActivationInterval: &utils.TPActivationInterval{
				ActivationTime: tps[0].ActivationInterval,
			},
			Stored:     tps[0].Stored,
			Blocker:    tps[0].Blocker,
			Weight:     tps[0].Weight,
			Limit:      tps[0].Limit,
			Thresholds: []string{"WARN_RES1", "WARN_RES2", "WARN3"},
		},
		&utils.TPResource{
			TPid: tps[2].Tpid,
			ID:   tps[2].Tag,
			Filters: []*utils.TPRequestFilter{
				&utils.TPRequestFilter{
					Type:      tps[2].FilterType,
					FieldName: tps[2].FilterFieldName,
					Values:    []string{"+40"},
				},
			},
			ActivationInterval: &utils.TPActivationInterval{
				ActivationTime: tps[2].ActivationInterval,
			},
			Stored:  tps[2].Stored,
			Blocker: tps[2].Blocker,
			Weight:  tps[2].Weight,
			Limit:   tps[2].Limit,
		},
	}
	rcvTPs := TpResources(tps).AsTPResources()
	if !(reflect.DeepEqual(eTPs, rcvTPs) || reflect.DeepEqual(eTPs[0], rcvTPs[1])) {
		t.Errorf("\nExpecting:\n%+v\nReceived:\n%+v", utils.ToIJSON(eTPs), utils.ToIJSON(rcvTPs))
	}
}

func TestAPItoResource(t *testing.T) {
	tpRL := &utils.TPResource{
		TPid: testTPID,
		ID:   "ResGroup1",
		Filters: []*utils.TPRequestFilter{
			&utils.TPRequestFilter{Type: MetaString, FieldName: "Account", Values: []string{"1001", "1002"}},
			&utils.TPRequestFilter{Type: MetaStringPrefix, FieldName: "Destination", Values: []string{"10", "20"}},
			&utils.TPRequestFilter{Type: MetaStatS, Values: []string{"CDRST1:*min_asr:34", "CDRST_1001:*min_asr:20"}},
			&utils.TPRequestFilter{Type: MetaRSRFields, Values: []string{"Subject(~^1.*1$)", "Destination(1002)"}},
		},
		ActivationInterval: &utils.TPActivationInterval{ActivationTime: "2014-07-29T15:00:00Z"},
		Stored:             false,
		Blocker:            false,
		Weight:             10,
		Limit:              "2",
	}
	eRL := &ResourceCfg{
		ID:      tpRL.ID,
		Stored:  tpRL.Stored,
		Blocker: tpRL.Blocker,
		Weight:  tpRL.Weight,
		Filters: make([]*RequestFilter, len(tpRL.Filters))}
	eRL.Filters[0] = &RequestFilter{Type: MetaString,
		FieldName: "Account", Values: []string{"1001", "1002"}}
	eRL.Filters[1] = &RequestFilter{Type: MetaStringPrefix,
		FieldName: "Destination", Values: []string{"10", "20"}}
	eRL.Filters[2] = &RequestFilter{Type: MetaStatS,
		Values: []string{"CDRST1:*min_asr:34", "CDRST_1001:*min_asr:20"},
		statSThresholds: []*RFStatSThreshold{
			&RFStatSThreshold{QueueID: "CDRST1", ThresholdType: "*min_asr", ThresholdValue: 34},
			&RFStatSThreshold{QueueID: "CDRST_1001", ThresholdType: "*min_asr", ThresholdValue: 20},
		}}
	eRL.Filters[3] = &RequestFilter{Type: MetaRSRFields, Values: []string{"Subject(~^1.*1$)", "Destination(1002)"},
		rsrFields: utils.ParseRSRFieldsMustCompile("Subject(~^1.*1$);Destination(1002)", utils.INFIELD_SEP),
	}
	at, _ := utils.ParseTimeDetectLayout("2014-07-29T15:00:00Z", "UTC")
	eRL.ActivationInterval = &utils.ActivationInterval{ActivationTime: at}
	eRL.Limit = 2
	if rl, err := APItoResource(tpRL, "UTC"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eRL, rl) {
		t.Errorf("Expecting: %+v, received: %+v", eRL, rl)
	}
}

func TestTPStatsAsTPStats(t *testing.T) {
	tps := []*TpStats{
		&TpStats{
			Tpid:               "TEST_TPID",
			Tag:                "Stats1",
			FilterType:         MetaStringPrefix,
			FilterFieldName:    "Account",
			FilterFieldValues:  "1001;1002",
			ActivationInterval: "2014-07-29T15:00:00Z",
			QueueLength:        100,
			TTL:                "1s",
			Metrics:            "*asr;*acd;*acc",
			Thresholds:         "THRESH1;THRESH2",
			Stored:             false,
			Blocker:            false,
			Weight:             20.0,
		},
	}
	eTPs := []*utils.TPStats{
		&utils.TPStats{
			TPid: tps[0].Tpid,
			ID:   tps[0].Tag,
			Filters: []*utils.TPRequestFilter{
				&utils.TPRequestFilter{
					Type:      tps[0].FilterType,
					FieldName: tps[0].FilterFieldName,
					Values:    []string{"1001", "1002"},
				},
			},
			ActivationInterval: &utils.TPActivationInterval{
				ActivationTime: tps[0].ActivationInterval,
			},
			QueueLength: tps[0].QueueLength,
			TTL:         tps[0].TTL,
			Metrics:     []string{"*asr", "*acd", "*acc"},
			Thresholds:  []string{"THRESH1", "THRESH2"},
			Stored:      tps[0].Stored,
			Blocker:     tps[0].Blocker,
			Weight:      tps[0].Weight,
		},
	}
	rcvTPs := TpStatsS(tps).AsTPStats()
	if !(reflect.DeepEqual(eTPs, rcvTPs) || reflect.DeepEqual(eTPs[0], rcvTPs[0])) {
		t.Errorf("\nExpecting:\n%+v\nReceived:\n%+v", utils.ToIJSON(eTPs), utils.ToIJSON(rcvTPs))
	}
}

func TestAPItoTPStats(t *testing.T) {
	tps := &utils.TPStats{
		TPid: testTPID,
		ID:   "Stats1",
		Filters: []*utils.TPRequestFilter{
			&utils.TPRequestFilter{Type: MetaString, FieldName: "Account", Values: []string{"1001", "1002"}},
		},
		ActivationInterval: &utils.TPActivationInterval{ActivationTime: "2014-07-29T15:00:00Z"},
		QueueLength:        100,
		TTL:                "1s",
		Metrics:            []string{"*asr", "*acd", "*acc"},
		Thresholds:         []string{"THRESH1", "THRESH2"},
		Stored:             false,
		Blocker:            false,
		Weight:             20.0,
	}

	eTPs := &StatsConfig{ID: tps.ID,
		QueueLength: tps.QueueLength,
		Metrics:     []string{"*asr", "*acd", "*acc"},
		Thresholds:  []string{"THRESH1", "THRESH2"},
		Filters:     make([]*RequestFilter, len(tps.Filters)),
		Stored:      tps.Stored,
		Blocker:     tps.Blocker,
		Weight:      20.0,
	}
	if eTPs.TTL, err = utils.ParseDurationWithSecs(tps.TTL); err != nil {
		t.Errorf("Got error: %+v", err)
	}

	eTPs.Filters[0] = &RequestFilter{Type: MetaString,
		FieldName: "Account", Values: []string{"1001", "1002"}}
	at, _ := utils.ParseTimeDetectLayout("2014-07-29T15:00:00Z", "UTC")
	eTPs.ActivationInterval = &utils.ActivationInterval{ActivationTime: at}

	if st, err := APItoStats(tps, "UTC"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eTPs, st) {
		t.Errorf("Expecting: %+v, received: %+v", eTPs, st)
	}
}

func TestAsTPThresholdAsAsTPThreshold(t *testing.T) {
	tps := []*TpThreshold{
		&TpThreshold{
			Tpid:               "TEST_TPID",
			Tag:                "Stats1",
			FilterType:         MetaStringPrefix,
			FilterFieldName:    "Account",
			FilterFieldValues:  "1001;1002",
			ActivationInterval: "2014-07-29T15:00:00Z",
			MinItems:           100,
			Recurrent:          false,
			MinSleep:           "1s",
			ThresholdType:      "",
			ThresholdValue:     1.2,
			Stored:             false,
			Blocker:            false,
			Weight:             20.0,
			ActionIDs:          "WARN3",
		},
	}
	eTPs := []*utils.TPThreshold{
		&utils.TPThreshold{
			TPid: tps[0].Tpid,
			ID:   tps[0].Tag,
			Filters: []*utils.TPRequestFilter{
				&utils.TPRequestFilter{
					Type:      tps[0].FilterType,
					FieldName: tps[0].FilterFieldName,
					Values:    []string{"1001", "1002"},
				},
			},
			ActivationInterval: &utils.TPActivationInterval{
				ActivationTime: tps[0].ActivationInterval,
			},
			MinItems:       tps[0].MinItems,
			MinSleep:       tps[0].MinSleep,
			ThresholdType:  tps[0].ThresholdType,
			ThresholdValue: tps[0].ThresholdValue,
			Recurrent:      tps[0].Recurrent,
			Stored:         tps[0].Stored,
			Blocker:        tps[0].Blocker,
			Weight:         tps[0].Weight,
			ActionIDs:      []string{"WARN3"},
		},
	}
	rcvTPs := TpThresholdS(tps).AsTPThreshold()
	if !(reflect.DeepEqual(eTPs, rcvTPs) || reflect.DeepEqual(eTPs[0], rcvTPs[0])) {
		t.Errorf("\nExpecting:\n%+v\nReceived:\n%+v", utils.ToIJSON(eTPs), utils.ToIJSON(rcvTPs))
	}
}

func TestAPItoTPThreshold(t *testing.T) {
	tps := &utils.TPThreshold{
		TPid: testTPID,
		ID:   "Stats1",
		Filters: []*utils.TPRequestFilter{
			&utils.TPRequestFilter{Type: MetaString, FieldName: "Account", Values: []string{"1001", "1002"}},
		},
		ActivationInterval: &utils.TPActivationInterval{ActivationTime: "2014-07-29T15:00:00Z"},
		MinItems:           100,
		Recurrent:          false,
		MinSleep:           "1s",
		ThresholdType:      "",
		ThresholdValue:     1.2,
		Stored:             false,
		Blocker:            false,
		Weight:             20.0,
		ActionIDs:          []string{"WARN3"},
	}

	eTPs := &ThresholdCfg{ID: tps.ID,
		Filters:        make([]*RequestFilter, len(tps.Filters)),
		MinItems:       tps.MinItems,
		Recurrent:      tps.Recurrent,
		ThresholdType:  tps.ThresholdType,
		ThresholdValue: tps.ThresholdValue,
		Stored:         tps.Stored,
		Blocker:        tps.Blocker,
		Weight:         tps.Weight,
		ActionIDs:      []string{"WARN3"},
	}
	if eTPs.MinSleep, err = utils.ParseDurationWithSecs(tps.MinSleep); err != nil {
		t.Errorf("Got error: %+v", err)
	}
	eTPs.Filters[0] = &RequestFilter{Type: MetaString,
		FieldName: "Account", Values: []string{"1001", "1002"}}
	at, _ := utils.ParseTimeDetectLayout("2014-07-29T15:00:00Z", "UTC")
	eTPs.ActivationInterval = &utils.ActivationInterval{ActivationTime: at}
	if st, err := APItoThresholdCfg(tps, "UTC"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(eTPs, st) {
		t.Errorf("Expecting: %+v, received: %+v", eTPs, st)
	}
}
