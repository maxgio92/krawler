package matrix

// The combinations are built concatenating one element per column into a row,
// and traversing all the columns' elements by shifting them from the last to the first columns
// (decremental abscissa order).
//
// (ordinate)
// y
// ^
// |	4
// |	3	Z
// |B	2	Y
// |A	1	X
// -----------> x (abscissa)
//
// Provided the sample scenario above, this function
// should combine the elements in the order below:
//
// A + 1 + X
// A + 1 + Y
// A + 1 + Z
// A + 2 + X
// ...
// B + 4 + Z
//
//nolint:godot
func GetColumnOrderedCombinationRows(columns []Column) ([]string, error) {
	rows := []string{}
	row := ""
	completed := false

	// For each time the last column has been reached
	// exit from recursion until reaching this:
	for {
		row = ""

		// Start always from the first column (x=0)
		err := gotoNextColumn(&rows, &row, 0, &columns[0], columns, &completed)
		if err != nil {
			return nil, err
		}

		ssp, ok := columns[0].Points.([]string)
		if !ok {
			return nil, NewErrUnsopportedPointType()
		}

		if columns[0].CurrentOrdinateIndex == len(ssp) || completed {
			break
		}
	}

	return rows, nil
}

func gotoNextColumn(points *[]string, row *string, abscissaIndex int, column *Column, columns []Column, completed *bool) error {
	currentColumnPoints, ok := column.Points.([]string)
	if !ok {
		return NewErrUnsopportedPointType()
	}

	if abscissaIndex+1 < len(columns) { // Until the last column is reached

		*row += currentColumnPoints[column.CurrentOrdinateIndex]

		// Move forward
		abscissaIndex++
		column = &columns[abscissaIndex]
		err := gotoNextColumn(points, row, abscissaIndex, column, columns, completed)
		if err != nil {
			return err
		}

	} else { // When the last column is reached

		for _, point := range currentColumnPoints {
			*points = append(*points, string(*row+point))
		}

		// Move backward
		if abscissaIndex > 0 {
			abscissaIndex--
		}
		column = &columns[abscissaIndex]

		// Store where we gone
		err := scrollDownPrevColumnPoint(column, columns, abscissaIndex, completed)
		if err != nil {
			return err
		}
	}
	return nil
}

func scrollDownPrevColumnPoint(column *Column, columns []Column, abscissaIndex int, completed *bool) error {
	currentColumnPoints, ok := column.Points.([]string)
	if !ok {
		return NewErrUnsopportedPointType()
	}

	// If the current column has still elements/points
	// and there are more than one column.
	if column.CurrentOrdinateIndex+1 < len(currentColumnPoints) && len(columns) > 1 {
		column.CurrentOrdinateIndex++

		// If the current column has been completely processed.
	} else {
		column.CurrentOrdinateIndex = 0
		abscissaIndex--

		// If it's not the first column.
		if abscissaIndex >= 0 {
			err := scrollDownPrevColumnPoint(&columns[abscissaIndex], columns, abscissaIndex, completed)
			if err != nil {
				return err
			}

			// If it's the first column.
		} else {
			*completed = true
		}
	}
	return nil
}
