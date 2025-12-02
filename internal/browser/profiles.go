package browser

var ChromeDesktopUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

func IsValidDeviceType(deviceType string) bool {
	return deviceType == "desktop" || deviceType == "mobile"
}

func IsValidDeviceSize(deviceSize string) bool {
	return deviceSize == "large" || deviceSize == "medium" || deviceSize == "small"
}

func DimensionsFromDeviceProfile(deviceType, deviceSize string) (h int, w int) {
	switch deviceType {
	case "desktop", "":
		switch deviceSize {
		case "large":
			return 1920, 1080
		case "medium", "":
			return 1536, 864
		case "small":
			return 1280, 720
		}
	case "mobile":
		switch deviceSize {
		case "large":
			return 430, 932
		case "medium", "":
			return 390, 844
		case "small":
			return 375, 812
		}
	}

	return 1280, 720
}

func UserAgent(deviceType, userAgentAlias string) string {
	switch userAgentAlias {
	case "chrome", "":
		switch deviceType {
		case "desktop", "":
			return ChromeDesktopUserAgent
		case "mobile":
			return "Mozilla/5.0 (iPhone; CPU iPhone OS 18_3_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3.1 Mobile/15E148 Safari/604"
		}
	}

	return ""
}
