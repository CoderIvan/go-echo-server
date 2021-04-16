package datagram

// Datagram *
type Datagram struct {
	TagName     string
	Addr        string
	ProjectName string
	Content     string
	ContextID   string
	ExtraInfo   string
	Time        int64
}

func (this Datagram) Equal(that Datagram) bool {
	if this.TagName == that.TagName {
		if this.Addr == that.Addr {
			if this.ProjectName == that.ProjectName {
				if this.ContextID == that.ContextID {
					if this.Content == that.Content {
						if this.ExtraInfo == that.ExtraInfo {
							if this.Time == that.Time {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}
