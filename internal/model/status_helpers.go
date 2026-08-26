package model

import "strings"

func (s SpringStatus) Valid() bool {
	switch s {
	case SpringActive, SpringInactive, SpringWatch:
		return true
	default:
		return false
	}
}

func (s SpringStatus) Label() string {
	switch s {
	case SpringActive:
		return "运行"
	case SpringInactive:
		return "停用"
	case SpringWatch:
		return "观察"
	default:
		return "未知"
	}
}

func (s SensorKind) Valid() bool {
	switch s {
	case SensorFlow, SensorTemperature, SensorConductivity, SensorLevel:
		return true
	default:
		return false
	}
}

func (s SensorKind) Label() string {
	switch s {
	case SensorFlow:
		return "流量"
	case SensorTemperature:
		return "温度"
	case SensorConductivity:
		return "电导率"
	case SensorLevel:
		return "水位"
	default:
		return "未知"
	}
}

func (s SensorStatus) Valid() bool {
	return s == SensorOnline || s == SensorOffline || s == SensorService
}

func (s SensorStatus) Label() string {
	switch s {
	case SensorOnline:
		return "在线"
	case SensorOffline:
		return "离线"
	case SensorService:
		return "维护中"
	default:
		return "未知"
	}
}

func (q ReadingQuality) Valid() bool {
	return q == QualityGood || q == QualitySuspect || q == QualityInvalid
}

func (q ReadingQuality) Weight() float64 {
	switch q {
	case QualityGood:
		return 1
	case QualitySuspect:
		return 0.5
	default:
		return 0
	}
}

func (p PulsePhase) Valid() bool {
	return p == PhaseRising || p == PhasePeaked || p == PhaseFading || p == PhaseConfirmed
}

func (p PulsePhase) Label() string {
	switch p {
	case PhaseRising:
		return "上升"
	case PhasePeaked:
		return "峰值"
	case PhaseFading:
		return "衰减"
	case PhaseConfirmed:
		return "确认"
	default:
		return "未知"
	}
}

func (s EventSeverity) Valid() bool {
	return s == SeverityInfo || s == SeverityWarning || s == SeverityCritical
}

func (s EventSeverity) Weight() int {
	switch s {
	case SeverityCritical:
		return 3
	case SeverityWarning:
		return 2
	default:
		return 1
	}
}

func (s EventSeverity) Label() string {
	switch s {
	case SeverityCritical:
		return "严重"
	case SeverityWarning:
		return "警告"
	default:
		return "提示"
	}
}

func (s BatchStatus) Valid() bool {
	return s == BatchOpen || s == BatchSubmitted || s == BatchArchived
}

func (s BatchStatus) Label() string {
	switch s {
	case BatchOpen:
		return "开放"
	case BatchSubmitted:
		return "已提交"
	case BatchArchived:
		return "已归档"
	default:
		return "未知"
	}
}

func (s AlertLevel) Valid() bool {
	return s == AlertInfo || s == AlertWarning || s == AlertCritical
}

func (s AlertStatus) Valid() bool {
	return s == AlertOpen || s == AlertAcknowledged
}

func (s AlertLevel) Weight() int {
	switch s {
	case AlertCritical:
		return 3
	case AlertWarning:
		return 2
	default:
		return 1
	}
}

func (s MaintenanceStatus) Valid() bool {
	return s == MaintenancePlanned || s == MaintenanceInProgress || s == MaintenanceDone || s == MaintenanceOverdue
}

func (s MaintenanceStatus) Label() string {
	switch s {
	case MaintenancePlanned:
		return "计划"
	case MaintenanceInProgress:
		return "进行中"
	case MaintenanceDone:
		return "完成"
	case MaintenanceOverdue:
		return "逾期"
	default:
		return "未知"
	}
}

func CanonicalStatus(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func IsOneOf(value string, allowed ...string) bool {
	canonical := CanonicalStatus(value)
	for _, candidate := range allowed {
		if canonical == CanonicalStatus(candidate) {
			return true
		}
	}
	return false
}
