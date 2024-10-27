package scenes

type Scene string

const (
	SceneDefault        Scene = "default"         // Состояние по умолчанию
	SceneEnterCity      Scene = "enter_city"      // Ввод города
	SceneSelectInterval Scene = "select_interval" // Выбор интервала
)
