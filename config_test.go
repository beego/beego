package beego

import (
	"testing"
)

func TestDefaults(t *testing.T) {
	if FlashName != "BEEGO_FLASH" {
		t.Errorf("FlashName was not set to default.")
	}

	if FlashSeperator != "BEEGOFLASH" {
		t.Errorf("FlashName was not set to default.")
	}
}
