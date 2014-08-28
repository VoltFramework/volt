package api

import "testing"

func TestSlice(t *testing.T) {
	var tasks = make(Tasks, 0)

	tasks = append(tasks, &Task{ID: "0"})
	tasks = append(tasks, &Task{ID: "1"})
	tasks = append(tasks, &Task{ID: "2"})
	tasks = append(tasks, &Task{ID: "3"})
	tasks = append(tasks, &Task{ID: "4"})

	if size := len(tasks.Slice(0, 0)); size != 0 {
		t.Fatalf("Expected size to be 0, got %d", size)
	}

	if size := len(tasks.Slice(0, 1)); size != 1 {
		t.Fatalf("Expected size to be 1, got %d", size)
	}

	if size := len(tasks.Slice(0, 20)); size != 5 {
		t.Fatalf("Expected size to be 5, got %d", size)
	}

	if size := len(tasks.Slice(1, 1)); size != 1 {
		t.Fatalf("Expected size to be 1, got %d", size)
	}

	if size := len(tasks.Slice(1, 20)); size != 0 {
		t.Fatalf("Expected size to be 0, got %d", size)
	}

	tasks = append(tasks, &Task{ID: "5"})
	tasks = append(tasks, &Task{ID: "6"})
	tasks = append(tasks, &Task{ID: "7"})
	tasks = append(tasks, &Task{ID: "8"})
	tasks = append(tasks, &Task{ID: "9"})
	tasks = append(tasks, &Task{ID: "10"})
	tasks = append(tasks, &Task{ID: "11"})
	tasks = append(tasks, &Task{ID: "12"})
	tasks = append(tasks, &Task{ID: "13"})
	tasks = append(tasks, &Task{ID: "14"})
	tasks = append(tasks, &Task{ID: "15"})
	tasks = append(tasks, &Task{ID: "16"})
	tasks = append(tasks, &Task{ID: "17"})
	tasks = append(tasks, &Task{ID: "18"})
	tasks = append(tasks, &Task{ID: "19"})
	tasks = append(tasks, &Task{ID: "20"})
	tasks = append(tasks, &Task{ID: "21"})
	tasks = append(tasks, &Task{ID: "22"})
	tasks = append(tasks, &Task{ID: "23"})
	tasks = append(tasks, &Task{ID: "24"})

	if size := len(tasks.Slice(0, 0)); size != 0 {
		t.Fatalf("Expected size to be 0, got %d", size)
	}

	if size := len(tasks.Slice(0, 1)); size != 1 {
		t.Fatalf("Expected size to be 1, got %d", size)
	}

	if size := len(tasks.Slice(0, 20)); size != 20 {
		t.Fatalf("Expected size to be 20, got %d", size)
	}

	if size := len(tasks.Slice(1, 1)); size != 1 {
		t.Fatalf("Expected size to be 1, got %d", size)
	}

	if size := len(tasks.Slice(1, 20)); size != 5 {
		t.Fatalf("Expected size to be 5, got %d", size)
	}
}
