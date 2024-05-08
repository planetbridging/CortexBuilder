package main

type EvaluationMetrics struct {
	TrainingError   float64
	ValidationError float64
	Accuracy        float64
	Precision       float64
	Recall          float64
	F1Score         float64
	AUCROC          float64
	MAE             float64
	MSE             float64
	RMSE            float64
}
