package store

import "testing"

func TestMemoryStore_CRUD(t *testing.T) {
	st := NewMemoryStore()

	created := st.Create("vm-1", "/subscriptions/x/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm-1", map[string]string{
		"env": "dev",
	})

	if created.ID == "" {
		t.Fatal("expected created resource to have an ID")
	}

	got, err := st.Get(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected same ID, got %q", got.ID)
	}

	all := st.List()
	if len(all) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(all))
	}

	if err := st.Delete(created.ID); err != nil {
		t.Fatalf("expected delete ok, got %v", err)
	}

	_, err = st.Get(created.ID)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestMemoryStore_Errors_TableDriven(t *testing.T) {
	st := NewMemoryStore()

	tests := []struct {
		name string
		run  func() error
		want error
	}{
		{
			name: "Get missing -> ErrNotFound",
			run: func() error {
				_, err := st.Get("missing")
				return err
			},
			want: ErrNotFound,
		},
		{
			name: "Delete missing -> ErrNotFound",
			run: func() error {
				return st.Delete("missing")
			},
			want: ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}
