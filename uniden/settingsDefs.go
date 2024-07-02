package uniden

import (
	"github.com/smoke7385/smk-uniden-bluetooth/types"
)

// SETTINGS DEFINITIONS
var defSettings = Settings{
	&Setting{
		Name: "Speed Cameras Alert Distance",
		StorageIndex: map[types.Model]int{
			types.R4: 8,
			types.R8: 9,
			types.R9: 11,
		},
		Values: Values{
			{"1000ft / 300m", 1},
			{"2000ft / 600m", 2},
			{"2500ft / 760m", 3},
			{"3000ft / 900m", 4},
			{"Auto", 5},
		},
	},
	&Setting{
		Name: "Enable Speed Cameras",
		StorageIndex: map[types.Model]int{
			types.R4: 7,
			types.R8: 8,
			types.R9: 10,
		},
		Values: BooleanValues,
	},
	&Setting{
		Name: "Alerts Priority",
		StorageIndex: map[types.Model]int{
			types.R4: 46,
			types.R8: 48,
			types.R9: 55,
		},
		Values: Values{
			{"SIGNAL", 0},
			{"KA_MRCD", 1},
			{"MRCD_KA", 2},
		},
	},
	&Setting{
		Name: "Auto mute memory option",
		StorageIndex: map[types.Model]int{
			types.R4: 95,
			types.R8: 51,
			types.R9: 58,
		},
		Values: BooleanValues,
	},
	&Setting{
		Name: "Enable Red Light Cameras",
		StorageIndex: map[types.Model]int{
			types.R4: 9,
			types.R8: 10,
			types.R9: 12,
		},
		Values: BooleanValues,
	},
	&Setting{
		Name: "Background Color",
		StorageIndex: map[types.Model]int{
			types.R4: 50,
			types.R8: 53,
			types.R9: 60,
		},
		Values: ColorValues,
	},
	&Setting{
		Name: "Quiet Ride Speed",
		StorageIndex: map[types.Model]int{
			types.R4: 77,
			types.R8: 87,
			types.R9: 104,
		},
		DynamicValues: func(s *Settings) *Values {
			return getSpeedValues(s, 5, 90, 5, 10, 90, 10)
		},
	},
	&Setting{
		Name: "Auto mute memory option",
		StorageIndex: map[types.Model]int{
			types.R4: 95,
			types.R8: 51,
			types.R9: 58,
		},
		Values: BooleanValues,
	},
	&Setting{
		Name: "Red light camera quiet ride speed",
		StorageIndex: map[types.Model]int{
			types.R4: 10,
			types.R8: 11,
			types.R9: 13,
		},
		DynamicValues: func(s *Settings) *Values {
			return getSpeedValues(s, 50, 85, 5, 80, 140, 10)
		},
	},
	&Setting{
		Name: "Operation mode",
		StorageIndex: map[types.Model]int{
			types.R4: 1,
			types.R8: 1,
			types.R9: 1,
		},
		Values: Values{
			{"Highway", 0},
			{"City", 1},
			{"Auto City", 2},
			{"Advanced", 3},
		},
	},
	&Setting{
		Name: "Auto City Mode Speed",
		StorageIndex: map[types.Model]int{
			types.R4: 5,
			types.R8: 5,
			types.R9: 8,
		},
		DynamicValues: func(s *Settings) *Values {
			return getSpeedValues(s, 10, 60, 5, 10, 100, 10)
		},
	},
	&Setting{
		Name: "Speed Units",
		StorageIndex: map[types.Model]int{
			types.R4: 60,
			types.R8: 68,
			types.R9: 86,
		},
		Values: Values{
			{"MPH", 0},
			{"KPH", 1},
		},
	},
	// BANDS - TODO: Add Directional Bands and Antenna toggle (R9 ONLY)
	&Setting{
		Name: "X Band",
		StorageIndex: map[types.Model]int{
			types.R4: 13,
			types.R8: 15,
			types.R9: 0,
		},
		Values:       BooleanValues,
		DefaultValue: false,
	},
	&Setting{
		Name: "K Band",
		StorageIndex: map[types.Model]int{
			types.R4: 14,
			types.R8: 16,
		},
		Values:       BooleanValues,
		DefaultValue: true,
	},
	&Setting{
		Name: "Ka Band",
		StorageIndex: map[types.Model]int{
			types.R4: 15,
			types.R8: 17,
		},
		Values:       BooleanValues,
		DefaultValue: true,
	},
	&Setting{
		Name: "Laser",
		StorageIndex: map[types.Model]int{
			types.R4: 16,
			types.R8: 18,
			types.R9: 25,
		},
		Values:       BooleanValues,
		DefaultValue: true,
	},
	&Setting{
		Name: "K POP",
		StorageIndex: map[types.Model]int{
			types.R4: 26,
			types.R8: 28,
			types.R9: 35,
		},
		Values:       BooleanValues,
		DefaultValue: false,
	},
	&Setting{
		Name: "Ka POP",
		StorageIndex: map[types.Model]int{
			types.R4: 29,
			types.R8: 31,
			types.R9: 38,
		},
		Values:       BooleanValues,
		DefaultValue: false,
	},
	// BAND SENSITIVITIES
	&Setting{
		Name: "X band sensitivity",
		StorageIndex: map[types.Model]int{
			types.R4: 2,
			types.R8: 2,
		},
		Values:       generateSlidersRange(30, 100, 10, "%"),
		DefaultValue: 100,
	},
	&Setting{
		Name: "K band sensitivity",
		StorageIndex: map[types.Model]int{
			types.R4: 3,
			types.R8: 3,
		},
		Values:       generateSlidersRange(30, 100, 10, "%"),
		DefaultValue: 100,
	},
	&Setting{
		Name: "Ka band sensitivity",
		StorageIndex: map[types.Model]int{
			types.R4: 4,
			types.R8: 4,
		},
		Values:       generateSlidersRange(30, 100, 10, "%"),
		DefaultValue: 100,
	},
	// BAND FILTERS
	&Setting{
		Name: "K band filter",
		StorageIndex: map[types.Model]int{
			types.R4: 30,
			types.R8: 32,
			types.R9: 39,
		},
		Values: BooleanValues,
	},
	&Setting{
		Name: "K block 24.199 (±0.002) filter",
		StorageIndex: map[types.Model]int{
			types.R4: 33,
			types.R8: 35,
		},
		Values: Values{
			Value{"OFF", 0},
			Value{"ON", 1},
			Value{"WEAK", 2},
		},
	},
	&Setting{
		Name: "K block 24.168 (±0.002) filter",
		StorageIndex: map[types.Model]int{
			types.R4: 34,
			types.R8: 36,
		},
		Values: Values{
			Value{"OFF", 0},
			Value{"ON", 1},
			Value{"WEAK", 2},
		},
	},
	// KA SCAN SEGMENTS
	&Setting{
		Name: "K scan width",
		StorageIndex: map[types.Model]int{
			types.R4: 35,
			types.R8: 37,
			types.R9: 44,
		},
		Values: Values{
			Value{"WIDE", 0},
			Value{"NARROW", 1},
			Value{"EXTENDED", 2},
		},
	},
	generateKaSegment(1, 37, 39, 46),
	generateKaSegment(2, 37, 39, 46),
	generateKaSegment(3, 37, 39, 46),
	generateKaSegment(4, 37, 39, 46),
	generateKaSegment(5, 37, 39, 46),
	generateKaSegment(6, 37, 39, 46),
	generateKaSegment(7, 37, 39, 46),
	generateKaSegment(8, 37, 39, 46),
	generateKaSegment(9, 37, 39, 46),

	// other 0

	&Setting{
		Name: "Auto mute volume",
		StorageIndex: map[types.Model]int{
			types.R4: 69,
			types.R8: 78,
			types.R9: 96,
		},
		Values: generateSlidersRange(0, 7, 1, ""),
	},

	&Setting{
		Name: "Auto mute memory option",
		StorageIndex: map[types.Model]int{
			types.R4: 95,
			types.R8: 51,
			types.R9: 58,
		},
		Values: Values{
			Value{"X_K", 0},
			Value{"X_K_KA", 1},
		},
	},
	&Setting{
		Name: "Mute memory option",
		StorageIndex: map[types.Model]int{
			types.R4: 47,
			types.R8: 49,
			types.R9: 56,
		},
		Values: Values{
			Value{"X_K", 0},
			Value{"X_K_KA", 1},
		},
	},

	&Setting{
		Name: "Quiet ride beep volume",
		StorageIndex: map[types.Model]int{
			types.R4: 79,
			types.R8: 89,
			types.R9: 106,
		},
		Values: generateSlidersRange(0, 8, 1, ""),
	},
	// Band Tones
	&Setting{
		Name: "Quiet ride beep volume",
		StorageIndex: map[types.Model]int{
			types.R4: 79,
			types.R8: 89,
			types.R9: 106,
		},
		Values: generateSlidersRange(0, 8, 1, ""),
	},
	&Setting{
		Name:   "X band tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 61,
			types.R8: 69,
			types.R9: 87,
		},
	},
	&Setting{
		Name:   "K band tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 62,
			types.R8: 70,
			types.R9: 88,
		},
	},
	&Setting{
		Name:   "Ka band tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 65,
			types.R8: 74,
			types.R9: 92,
		},
	},
	&Setting{
		Name:   "MRCD/T tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 63,
			types.R8: 72,
			types.R9: 90,
		},
	},
	&Setting{
		Name:   "Gatso tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 64,
			types.R8: 73,
			types.R9: 91,
		},
	},
	&Setting{
		Name:   "Laser tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 67,
			types.R8: 76,
			types.R9: 94,
		},
	},
	&Setting{
		Name:   "K band bogey tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 93,
			types.R8: 71,
			types.R9: 89,
		},
	},
	&Setting{
		Name:   "Ka band bogey tone",
		Values: ToneValues,
		StorageIndex: map[types.Model]int{
			types.R4: 66,
			types.R8: 75,
			types.R9: 93,
		},
	},
	&Setting{
		Name: "Alerts priority",
		StorageIndex: map[types.Model]int{
			types.R4: 46,
			types.R8: 48,
			types.R9: 55,
		},
		Values: Values{
			Value{"SIGNAL", 0},
			Value{"KA_MRCD", 1},
			Value{"MRCD_KA", 2},
		},
	},
	&Setting{
		Name: "Limit speed",
		StorageIndex: map[types.Model]int{
			types.R4: 80,
			types.R8: 90,
			types.R9: 107,
		},
		DynamicValues: func(s *Settings) *Values {
			values := append(Values{
				Value{"Off", 0},
			}, *getSpeedValuesLiteral(s, 50, 100, 5, 80, 160, 10)...)

			return &values
		},
	},
	&Setting{
		Name: "Display mode",
		StorageIndex: map[types.Model]int{
			types.R4: 56,
			types.R8: 64,
			types.R9: 83,
		},
		Values: Values{
			Value{"SCAN", 0},
			Value{"MODE", 1},
			Value{"TIME", 2},
		},
	},
	&Setting{
		Name: "Alert dsplay mode",
		StorageIndex: map[types.Model]int{
			types.R4: 59,
			types.R8: 67,
			types.R9: 155,
		},
		Values: Values{
			Value{"DISPLAY_1", 0},
			Value{"DISPLAY_2", 1},
			Value{"DISPLAY_3", 2},
		},
	},
	&Setting{
		Name: "Left display",
		StorageIndex: map[types.Model]int{
			types.R4: 58,
			types.R8: 66,
			types.R9: 85,
		},
		Values: Values{
			Value{"SPEED", 0},
			Value{"SPEED_COMPASS", 1},
			Value{"COMPASS", 2},
			Value{"VOLTAGE", 3},
			Value{"ALTITUDE", 4},
		},
	},
	&Setting{
		Name: "Left display",
		StorageIndex: map[types.Model]int{
			types.R4: 58,
			types.R8: 66,
			types.R9: 85,
		},
		Values: Values{
			Value{"SPEED", 0},
			Value{"SPEED_COMPASS", 1},
			Value{"COMPASS", 2},
			Value{"VOLTAGE", 3},
			Value{"ALTITUDE", 4},
		},
	},
	&Setting{
		Name: "X band color",
		StorageIndex: map[types.Model]int{
			types.R4: 51,
			types.R8: 59,
		},
		Values: BandColorValues,
	},
	&Setting{
		Name: "K band color",
		StorageIndex: map[types.Model]int{
			types.R4: 52,
			types.R8: 60,
		},
		Values: BandColorValues,
	},
	&Setting{
		Name: "Ka band color",
		StorageIndex: map[types.Model]int{
			types.R4: 55,
			types.R8: 53,
		},
		Values: BandColorValues,
	},
	&Setting{
		Name: "MRCD/T color",
		StorageIndex: map[types.Model]int{
			types.R4: 53,
			types.R8: 61,
		},
		Values: BandColorValues,
	},
	&Setting{
		Name: "Gatso color",
		StorageIndex: map[types.Model]int{
			types.R4: 54,
			types.R8: 62,
		},
		Values: BandColorValues,
	},
	&Setting{
		Name: "Display brightness",
		StorageIndex: map[types.Model]int{
			types.R4: 92,
			types.R8: 102,
			types.R9: 119,
		},
		Values: Values{
			{"OFF", 0},
			{"DARK", 1},
			{"DIMMER", 2},
			{"DIM", 3},
			{"BRIGHT", 4},
			{"AUTO", 5},
		},
	},
	&Setting{
		Name: "Dark mode",
		StorageIndex: map[types.Model]int{
			types.R4: 70,
			types.R8: 80,
			types.R9: 97,
		},
		Values: Values{
			{"DIMMER", 0},
			{"DIM", 1},
			{"BRIGHT", 2},
		},
	},
	&Setting{
		Name: "Bright brightness",
		StorageIndex: map[types.Model]int{
			types.R4: 73,
			types.R8: 83,
			types.R9: 100,
		},
		Values: Values{
			{"DIMMER", 0},
			{"DIM", 1},
			{"BRIGHT", 2},
		},
	},
	&Setting{
		Name: "Dim brightness",
		StorageIndex: map[types.Model]int{
			types.R4: 75,
			types.R8: 85,
			types.R9: 102,
		},
		Values: Values{
			{"OFF", 0},
			{"DARK", 1},
			{"DIMMER", 2},
			{"DIM", 3},
			{"BRIGHT", 4},
		},
	},
	&Setting{
		Name: "Auto dim mode",
		StorageIndex: map[types.Model]int{
			types.R4: 70,
			types.R8: 80,
			types.R9: 97,
		},
		Values: Values{
			{"SENSOR", 0},
			{"TIME", 1},
		},
	},
	&Setting{
		Name: "Bright time",
		StorageIndex: map[types.Model]int{
			types.R4: 72,
			types.R8: 82,
			types.R9: 99,
		},
		Values: Values{
			{"T_5_30", 0},
			{"T_5_45", 1},
			{"T_6_00", 2},
			{"T_6_15", 3},
			{"T_6_30", 4},
			{"T_6_45", 5},
			{"T_7_00", 6},
			{"T_7_15", 7},
			{"T_7_30", 8},
		},
	},
	&Setting{
		Name: "Dim time",
		StorageIndex: map[types.Model]int{
			types.R4: 74,
			types.R8: 84,
			types.R9: 101,
		},
		Values: Values{
			{"T_5_00", 0},
			{"T_5_15", 1},
			{"T_5_30", 2},
			{"T_5_45", 3},
			{"T_6_00", 4},
			{"T_6_15", 5},
			{"T_6_30", 6},
			{"T_6_45", 7},
			{"T_7_00", 8},
			{"T_7_15", 9},
			{"T_7_30", 10},
			{"T_7_45", 11},
			{"T_8_00", 12},
		},
	},
	&Setting{
		Name: "Time zone",
		StorageIndex: map[types.Model]int{
			types.R4: 81,
			types.R8: 91,
			types.R9: 108,
		},
		Values: Values{
			{"GMT-12", 0},
			{"GMT-11", 1},
			{"GMT-10", 2},
			{"GMT-9", 3},
			{"GMT-8", 4},
			{"GMT-7", 5},
			{"GMT-6", 6},
			{"GMT-5", 7},
			{"GMT-4", 8},
			{"GMT-3", 9},
			{"GMT-2", 10},
			{"GMT-1", 11},
			{"GMT", 12},
			{"GMT+1", 13},
			{"GMT+2", 14},
			{"GMT+3", 15},
			{"GMT+4", 16},
			{"GMT+5", 17},
			{"GMT+6", 18},
			{"GMT+7", 19},
			{"GMT+8", 20},
			{"GMT+9", 21},
			{"GMT+10", 22},
			{"GMT+11", 23},
			{"GMT+12", 24},
		},
	},

	&Setting{
		Name: "Detector volume",
		StorageIndex: map[types.Model]int{
			types.R4: 91,
			types.R8: 101,
			types.R9: 118,
		},
		Values: Values{
			{"Always Muted", 0},
			{"1", 1},
			{"2", 2},
			{"3", 3},
			{"4", 4},
			{"5", 5},
			{"6", 6},
			{"7", 7},
			{"8", 8},
		},
	},
	&Setting{
		Name: "Memory Quota",
		StorageIndex: map[types.Model]int{
			types.R4: 90,
			types.R8: 100,
			types.R9: 117,
		},
		Values: Values{
			{"UM_MM_1750_250", 0},
			{"UM_MM_1700_300", 1},
			{"UM_MM_1650_350", 2},
			{"UM_MM_1600_400", 3},
			{"UM_MM_1550_450", 4},
			{"UM_MM_1500_500", 5},
			{"UM_MM_1450_550", 6},
			{"UM_MM_1400_600", 7},
			{"UM_MM_1350_650", 8},
			{"UM_MM_1300_700", 9},
			{"UM_MM_1250_750", 10},
			{"UM_MM_1200_800", 11},
			{"UM_MM_1150_850", 12},
			{"UM_MM_1100_900", 13},
			{"UM_MM_1050_950", 14},
			{"UM_MM_1000_1000", 15},
			{"UM_MM_950_1050", 16},
			{"UM_MM_900_1100", 17},
			{"UM_MM_850_1150", 18},
			{"UM_MM_800_1200", 19},
			{"UM_MM_750_1250", 20},
			{"UM_MM_700_1300", 21},
			{"UM_MM_650_1350", 22},
			{"UM_MM_600_1400", 23},
			{"UM_MM_550_1450", 24},
			{"UM_MM_500_1500", 25},
			{"UM_MM_450_1550", 26},
			{"UM_MM_400_1600", 27},
			{"UM_MM_350_1650", 28},
			{"UM_MM_300_1700", 29},
			{"UM_MM_250_1750", 30},
		},
	},
	generateBool(78, 88, 105)("Enable quiet ride for MRCD/T"),
	generateBool(82, 92, 109)("Daylight Savings Time (DST)"),
	generateBool(83, 93, 110)("Low battery voltage warning"),
	generateBool(48, 50, 57)("Enable auto mute memory"),
	generateBool(84, 94, 111)("Vehicle battery saver"),
	generateBool(57, 65, 84)("All threat display"),
	generateBool(12, 14, 16)("KA frequency voice"),
	generateBool(68, 77, 95)("Enable auto mute"),
	generateBool(31, 33, 40)("Ka band filter"),
	generateBool(28, 30, 37)("Ka band filter"),
	generateBool(49, 12, 14)("POI Passchime"),
	generateBool(17, 19, 26)("Laser gun ID"),
	generateBool(11, 13, 15)("Enable voice"),
	generateBool(85, 95, 112)("Self test"),
	generateBool(76, 86, 103)("Backlight"),
	generateBool(57, 65, 84)("Scan icon"),
	generateBool(27, 29, 36)("MRCD/T"),
	generateBool(32, 34, 41)("TSF"),
	generateBool(6, 7, 9)("GPS"),
}
