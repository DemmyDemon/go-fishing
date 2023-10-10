package terrain

var Sand = Terrain{
	DecorationSpacing:         180,
	DecorationSpacingVariance: 150,
	DecorationScaleX:          0.5,
	DecorationScaleY:          1.0,
	TileSize:                  128,
	TilePaths: []string{
		"img/terrain/sand_05.png",
		"img/terrain/sand_06.png",
		"img/terrain/sand_07.png",
		"img/terrain/sand_08.png",
	},
	DecorationPaths: []string{
		"img/decoration/green_kelp_large.png",
		"img/decoration/green_kelp_small.png",
		"img/decoration/green_seagrass_large.png",
		"img/decoration/green_seagrass_small.png",
		"img/decoration/green_seaweed_large.png",
		"img/decoration/green_seaweed_small.png",
		"img/decoration/orange_coral_large.png",
		"img/decoration/orange_coral_small.png",
		"img/decoration/purple_kelp_large.png",
		"img/decoration/purple_kelp_small.png",
		"img/decoration/purple_seaweed_large.png",
		"img/decoration/purple_seaweed_small.png",
		"img/decoration/rock_round.png",
		"img/decoration/rock_square.png",
	},
}
