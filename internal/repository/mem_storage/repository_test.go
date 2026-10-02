package mem_storage

//func TestAddAndGet(t *testing.T) {
//	storage := NewMemStorage()
//	testVal := 1.1
//	tableTest := []struct {
//		name     string
//		metric   model.Metrics
//		expected model.Metrics
//		err      error
//	}{
//		{
//			"add and check",
//			model.Metrics{
//				ID:    "Test",
//				MType: model.Gauge,
//				Delta: nil,
//				Value: &testVal,
//				Hash:  "",
//			},
//			model.Metrics{
//				ID:    "Test",
//				MType: model.Gauge,
//				Delta: nil,
//				Value: &testVal,
//				Hash:  "",
//			},
//			nil,
//		},
//	}
//	for _, tt := range tableTest {
//		t.Run(tt.name, func(t *testing.T) {
//			err := storage.Add(t.Context(), tt.metric)
//			if err != nil {
//				assert.Equal(t, tt.err, err)
//				return
//			}
//			metric, err := storage.GetValue(t.Context(), tt.metric.ID)
//			assert.Equal(t, tt.expected, metric)
//		})
//	}
//}
