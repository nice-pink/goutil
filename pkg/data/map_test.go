package data

import "testing"

func TestPatchMapOverwrite(t *testing.T) {
}

func TestPatchMap(t *testing.T) {
	// input
	in := map[string]any{}
	sub1 := map[string]any{"key1": "val1"}
	sub2 := map[string]any{"key2": "val2"}
	in["sub1"] = sub1
	in["sub2"] = sub2

	// patch
	patch := map[string]any{}
	pSub2 := map[string]any{"key2": "val3"}
	patch["sub2"] = pSub2

	// patch sub2
	out := PatchMap(in, patch)
	vSub2 := out["sub2"].(map[string]any)
	if vSub2["key2"] != "val3" {
		t.Error("TestPatchMap:: patch sub2 failed. val3 ==", vSub2["key2"])
	}

	// patch sub4
	pSub4 := map[string]any{"key4": "val4"}
	patch["sub4"] = pSub4
	out = PatchMap(in, patch)
	vSub2 = out["sub2"].(map[string]any)
	if vSub2["key2"] != "val3" {
		t.Error("TestPatchMap:: patch sub3 failed. val3 ==", vSub2["key2"])
	}
	vSub4 := out["sub4"].(map[string]any)
	if vSub4["key4"] != "val4" {
		t.Error("TestPatchMap:: patch sub3 failed. val4 ==", vSub4["key4"])
	}
}
