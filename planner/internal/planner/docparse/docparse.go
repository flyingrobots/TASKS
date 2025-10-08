package docparse

import (
    "bufio"
    "fmt"
    "regexp"
    "strings"
    "encoding/json"

    m "github.com/james/tasks-planner/internal/model"
)

// Feature is a minimal representation parsed from a spec document.
type Feature struct {
    ID    string
    Title string
}

// TaskSpec is a minimal parsed task definition.
type TaskSpec struct {
    FeatureID string
    Title     string
    After     []string  // dependencies by title or ID
    Hours     float64   // duration hint in hours (0 if unset)
    Accept    []m.AcceptanceCheck
}

var (
    reFeature = regexp.MustCompile(`^\s{0,3}#{2}\s+(.+?)\s*$`) // lines starting with '## '
    reTask    = regexp.MustCompile(`^\s*[-*]\s+(?:\[.?\]\s*)?(.+?)\s*$`) // '- task title' or '- [ ] task'
    reAfter   = regexp.MustCompile(`(?i)\bafter\s*:\s*([^;]+)$`)             // 'after: A, B, T001'
    reDur     = regexp.MustCompile(`\((\d+(?:\.\d+)?)(h|m)\)`)             // '(3h)' or '(90m)'
)

// ParseMarkdown extracts features (## headings) and tasks (bullet items under last feature).
// It is intentionally simple: one level of features; tasks inherit the most recent feature.
func ParseMarkdown(input string) (features []Feature, tasks []TaskSpec) {
    scanner := bufio.NewScanner(strings.NewReader(input))
    scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

    var currentFeatureID string
    featureCount := 0
    lastTaskIdx := -1
    inFence := false
    fenceLang := ""
    var fenceBuf []string
    for scanner.Scan() {
        line := scanner.Text()
        // Fenced code blocks for acceptance
        if strings.HasPrefix(line, "```") {
            lang := strings.TrimSpace(strings.TrimPrefix(line, "```"))
            if !inFence {
                inFence = true
                fenceLang = strings.ToLower(lang)
                fenceBuf = fenceBuf[:0]
                continue
            } else {
                // closing fence
                if inFence && (fenceLang == "accept" || fenceLang == "acceptance" || fenceLang == "checks") && lastTaskIdx >= 0 {
                    payload := strings.Join(fenceBuf, "\n")
                    var arr []m.AcceptanceCheck
                    if strings.HasPrefix(strings.TrimSpace(payload), "[") {
                        _ = json.Unmarshal([]byte(payload), &arr)
                    } else {
                        var one m.AcceptanceCheck
                        if err := json.Unmarshal([]byte(payload), &one); err == nil { arr = []m.AcceptanceCheck{one} }
                    }
                    if len(arr) > 0 {
                        tasks[lastTaskIdx].Accept = append(tasks[lastTaskIdx].Accept, arr...)
                    }
                }
                inFence = false
                fenceLang = ""
                fenceBuf = fenceBuf[:0]
                continue
            }
        }
        if inFence {
            fenceBuf = append(fenceBuf, line)
            continue
        }
        if m := reFeature.FindStringSubmatch(line); m != nil {
            featureCount++
            id := formatID("F", featureCount)
            features = append(features, Feature{ID: id, Title: strings.TrimSpace(m[1])})
            currentFeatureID = id
            continue
        }
        if m := reTask.FindStringSubmatch(line); m != nil {
            raw := strings.TrimSpace(m[1])
            title := raw
            // duration
            hrs := 0.0
            if dm := reDur.FindStringSubmatch(raw); dm != nil {
                val := dm[1]
                unit := dm[2]
                if f, err := parseNumber(val); err == nil {
                    if unit == "h" { hrs = f } else { hrs = f / 60.0 }
                }
                title = strings.TrimSpace(strings.Replace(raw, dm[0], "", 1))
            }
            // after: ...
            after := []string{}
            if am := reAfter.FindStringSubmatch(title); am != nil {
                rhs := am[1]
                parts := strings.Split(rhs, ",")
                for _, p := range parts {
                    s := strings.TrimSpace(p)
                    if s != "" { after = append(after, s) }
                }
                // strip the 'after: ...' tail
                title = strings.TrimSpace(strings.Replace(title, am[0], "", 1))
            }
            if title == "" { continue }
            if currentFeatureID == "" {
                // create a default feature if none seen yet
                featureCount++
                currentFeatureID = formatID("F", featureCount)
                features = append(features, Feature{ID: currentFeatureID, Title: "General"})
            }
            tasks = append(tasks, TaskSpec{FeatureID: currentFeatureID, Title: title, After: after, Hours: hrs})
            lastTaskIdx = len(tasks) - 1
        }
    }
    return features, tasks
}

func formatID(prefix string, n int) string {
    return fmt.Sprintf("%s%03d", prefix, n)
}

func parseNumber(s string) (float64, error) {
    // Minimal float parser to avoid importing strconv changes later
    // Use fmt.Sscan for simplicity here
    var f float64
    _, err := fmt.Sscan(s, &f)
    return f, err
}
