package main

import "errors"

type stage struct {
	level         int
	mapPath       string
	hunterBuilder iEnemyBuilder
	ghostBuilder  iEnemyBuilder
	enemies       []iEnemy
}

func initStages() []stage {
	defaultHunterBuilder := newEnemyBuilder().defaultHunter()
	defaultGhostBuilder := newEnemyBuilder().defaultGhost()
	return []stage{
		{
			level:         1,
			mapPath:       "files/stage/map01.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
		},
		{
			level:         2,
			mapPath:       "files/stage/map02.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
		},
	}
}

func getStage(stages []stage, level int) (stage, error) {
	for _, stage := range stages {
		if level == stage.level {
			return stage, nil
		}
	}
	return stage{}, errors.New("File does not exist: " + stages[level].mapPath)
}
