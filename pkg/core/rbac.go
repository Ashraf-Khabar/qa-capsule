package core

// Global role identifiers (stored in users.role). UI: Platform Admin, Manager, Lead, Observer.
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleLead     = "lead"
	RoleObserver = "observer"
)

// Legacy codes migrated on startup (operator → lead, viewer → observer).
const (
	RoleOperatorLegacy = "operator"
	RoleViewerLegacy   = "viewer"
)

// NormalizeRole maps deprecated role codes to current ones.
func NormalizeRole(role string) string {
	switch role {
	case RoleOperatorLegacy:
		return RoleLead
	case RoleViewerLegacy:
		return RoleObserver
	default:
		return role
	}
}

// RoleLevel returns a numeric rank for hierarchy comparisons.
func RoleLevel(role string) int {
	switch NormalizeRole(role) {
	case RoleAdmin:
		return 4
	case RoleManager:
		return 3
	case RoleLead:
		return 2
	case RoleObserver:
		return 1
	default:
		return 0
	}
}

// HasMinRole is true when userRole meets or exceeds minRole in the hierarchy.
func HasMinRole(userRole, minRole string) bool {
	return RoleLevel(userRole) >= RoleLevel(minRole)
}

// IsValidRole reports whether role is one of the supported global roles (legacy codes accepted).
func IsValidRole(role string) bool {
	switch NormalizeRole(role) {
	case RoleAdmin, RoleManager, RoleLead, RoleObserver:
		return true
	default:
		return false
	}
}

// IsCanonicalRole is true only for current stored codes (no legacy operator/viewer).
func IsCanonicalRole(role string) bool {
	switch role {
	case RoleAdmin, RoleManager, RoleLead, RoleObserver:
		return true
	default:
		return false
	}
}

func IsManager(role string) bool {
	return NormalizeRole(role) == RoleManager
}

func IsAdmin(role string) bool {
	return NormalizeRole(role) == RoleAdmin
}

func IsLead(role string) bool {
	return NormalizeRole(role) == RoleLead
}

func IsObserver(role string) bool {
	return NormalizeRole(role) == RoleObserver
}

func CanAccessFinOps(role string) bool {
	return IsManager(role)
}

func CanManageTeams(role string) bool {
	return IsAdmin(role) || IsManager(role)
}

func CanManageIAM(role string) bool {
	return IsAdmin(role)
}

func CanManageProjects(role string) bool {
	return IsManager(role) || IsLead(role)
}

func CanAccessPlugins(role string) bool {
	return IsManager(role) || IsLead(role)
}

// CanManagePluginAutoRun allows toggling AUTO-RUN / gateway routing on integrations (Manager or Platform Admin).
func CanManagePluginAutoRun(role string) bool {
	return IsManager(role) || IsAdmin(role)
}

func CanResolveIncidents(role string) bool {
	return HasMinRole(role, RoleLead)
}

func CanDeleteIncidents(role string) bool {
	r := NormalizeRole(role)
	return r == RoleManager || r == RoleLead
}

func CanAccessChartStudio(role string) bool {
	r := NormalizeRole(role)
	return r == RoleManager || r == RoleLead || r == RoleObserver
}

// CanManageWorkflow allows creating/editing visual remediation DAGs (Lead or Manager).
func CanManageWorkflow(role string) bool {
	return IsManager(role) || IsLead(role)
}

// CanViewHealing allows reading self-healing context and insights (Observer+ on operations roles).
func CanViewHealing(role string) bool {
	r := NormalizeRole(role)
	return r == RoleLead || r == RoleManager || r == RoleObserver
}

// CanViewQuarantine allows reading quarantine list.
func CanViewQuarantine(role string) bool {
	return CanViewHealing(role)
}

// CanManageQuarantine allows manual quarantine / lift (Lead+).
func CanManageQuarantine(role string) bool {
	return HasMinRole(role, RoleLead)
}

// CanPatchExecutionFlags allows tagging pipeline runs (Lead+).
func CanPatchExecutionFlags(role string) bool {
	return HasMinRole(role, RoleLead)
}

// CanAccessOperations is the telemetry dashboard and pipeline reports. Platform Admin uses Workspaces instead.
func CanAccessOperations(role string) bool {
	return NormalizeRole(role) != RoleAdmin
}

// CanViewAIUsage is the AI & MCP guide (operations roles).
func CanViewAIUsage(role string) bool {
	r := NormalizeRole(role)
	return r == RoleLead || r == RoleManager || r == RoleObserver
}

// CanAccessDORA is the executive DORA dashboard (Manager only).
func CanAccessDORA(role string) bool {
	return IsManager(role)
}

// CanAccessView mirrors the web navigation guard for encrypted page links.
func CanAccessView(role, view string) bool {
	switch view {
	case "dashboard", "unified-report":
		return CanAccessOperations(role)
	case "about", "profile":
		return true
	case "finops":
		return CanAccessFinOps(role)
	case "organizations":
		return CanManageTeams(role)
	case "management", "settings":
		return CanManageIAM(role)
	case "plugins":
		return CanAccessPlugins(role)
	case "ingestion":
		return CanManageProjects(role)
	case "healing", "rca":
		return CanViewHealing(role)
	case "ai-config":
		return CanConfigureAI(role)
	case "ai-usage":
		return CanViewAIUsage(role)
	case "dora":
		return CanAccessDORA(role)
	default:
		return false
	}
}

// CanConfigureAI allows changing LLM provider settings (Manager only).
func CanConfigureAI(role string) bool {
	return NormalizeRole(role) == RoleManager
}

// CanManageHealing allows proposing fixes and opening remediation PRs (Lead+).
func CanManageHealing(role string) bool {
	return HasMinRole(role, RoleLead)
}
