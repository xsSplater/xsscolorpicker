// xsscolorpicker/style.go

package colorpicker

import "image/color"

// Style определяет цвета служебных элементов пикера:
//   - маркеров (курсоров на area и на полосах);
//   - шахматного фона под alpha-полосой.
//
// Основные цвета пикера (градиенты) строятся из самого цвета, поэтому
// в Style их нет.
type Style struct {
	// MarkerFill — заливка маркера.
	MarkerFill color.Color
	// MarkerStroke — обводка маркера. Должна контрастировать и с
	// градиентом, и с заливкой.
	MarkerStroke color.Color
	// CheckerLight / CheckerDark — два цвета шахматного узора.
	// Классический вариант — оттенки серого разной светлоты.
	CheckerLight color.Color
	CheckerDark  color.Color
	// CheckerBoxSize — размер одной клетки шахматки в пикселях.
	// Значение <= 0 трактуется как 10.
	CheckerBoxSize int
}

// DefaultStyle возвращает нейтрально-серый стиль. Используется, если
// приложение не вызвало SetDefaultStyle.
func DefaultStyle() Style {
	return Style{
		MarkerFill:     color.NRGBA{R: 50, G: 50, B: 50, A: 120},
		MarkerStroke:   color.NRGBA{R: 50, G: 50, B: 50, A: 200},
		CheckerLight:   color.Gray{Y: 58},
		CheckerDark:    color.Gray{Y: 84},
		CheckerBoxSize: 10,
	}
}

// currentStyle — стиль, применяемый ко всем вновь созданным пикерам.
// Уже созданные пикеры текущий стиль не перечитывают: при смене темы
// надо пересоздать окно, где живёт пикер.
//
// Так сделано намеренно — не расширяем интерфейс ColorPicker и не
// трогаем четыре реализации ради одной опции.
var currentStyle = DefaultStyle()

// SetDefaultStyle задаёт стиль для всех новых пикеров.
// Вызывать из главного потока Fyne, до создания пикера.
func SetDefaultStyle(s Style) {
	currentStyle = s
}

// CurrentStyle возвращает текущий глобальный стиль.
func CurrentStyle() Style {
	return currentStyle
}
