package resources

import (
	"os"
	"testing"
)

// TestResourceManagerCaching проверяет, что LoadImage дважды возвращает тот же указатель
func TestResourceManagerCaching(t *testing.T) {
	rm := NewResourceManager()

	// Используем существующую тестовую текстуру
	path := "internal/asset/image/dog1.png"

	// Загружаем изображение первый раз
	img1, err := rm.LoadImage(path)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed, got error: %v", err)
	}
	if img1 == nil {
		t.Fatal("Expected LoadImage to return non-nil image")
	}

	// Загружаем то же изображение второй раз
	img2, err := rm.LoadImage(path)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed on second call, got error: %v", err)
	}
	if img2 == nil {
		t.Fatal("Expected LoadImage to return non-nil image on second call")
	}

	// Проверяем, что это тот же указатель (кэширование работает)
	if img1 != img2 {
		t.Error("Expected LoadImage to return the same image instance (cached), but got different instances")
	}
}

// TestResourceManagerLoadError проверяет обработку ошибок загрузки
func TestResourceManagerLoadError(t *testing.T) {
	rm := NewResourceManager()

	// Пытаемся загрузить несуществующий файл
	nonExistentPath := "internal/asset/image/nonexistent_file_12345.png"

	img, err := rm.LoadImage(nonExistentPath)

	// Проверяем, что возвращается ошибка
	if err == nil {
		t.Error("Expected LoadImage to return error for non-existent file, got nil")
	}

	// Проверяем, что изображение nil
	if img != nil {
		t.Error("Expected LoadImage to return nil image on error, got non-nil")
	}

	// Проверяем, что ошибка связана с отсутствием файла
	if !os.IsNotExist(err) {
		t.Errorf("Expected error to be file not found error, got: %v", err)
	}
}

// TestResourceManagerMultipleImages проверяет кэширование нескольких изображений
func TestResourceManagerMultipleImages(t *testing.T) {
	rm := NewResourceManager()

	path1 := "internal/asset/image/dog1.png"
	path2 := "internal/asset/image/dog_jump.png"

	// Загружаем первое изображение
	img1a, err := rm.LoadImage(path1)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed for path1, got error: %v", err)
	}

	// Загружаем второе изображение
	img2a, err := rm.LoadImage(path2)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed for path2, got error: %v", err)
	}

	// Загружаем первое изображение снова
	img1b, err := rm.LoadImage(path1)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed for path1 (second time), got error: %v", err)
	}

	// Загружаем второе изображение снова
	img2b, err := rm.LoadImage(path2)
	if err != nil {
		t.Fatalf("Expected LoadImage to succeed for path2 (second time), got error: %v", err)
	}

	// Проверяем, что каждое изображение кэшируется отдельно
	if img1a != img1b {
		t.Error("Expected path1 to return same cached instance")
	}
	if img2a != img2b {
		t.Error("Expected path2 to return same cached instance")
	}

	// Проверяем, что разные пути возвращают разные изображения
	if img1a == img2a {
		t.Error("Expected different paths to return different images")
	}
}

// TestResourceManagerEmptyPath проверяет обработку пустого пути
func TestResourceManagerEmptyPath(t *testing.T) {
	rm := NewResourceManager()

	img, err := rm.LoadImage("")

	// Проверяем, что возвращается ошибка
	if err == nil {
		t.Error("Expected LoadImage to return error for empty path, got nil")
	}

	// Проверяем, что изображение nil
	if img != nil {
		t.Error("Expected LoadImage to return nil image for empty path, got non-nil")
	}
}
