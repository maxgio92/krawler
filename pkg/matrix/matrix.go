package matrix

// (ordinate)
// y
// ^
// |		4
// |		3	Z
// |	B	2	Y
// |	A	1	X
// ---------------> x (abscissa)

//nolint:godot
// Provided the sample scenario above, this function
// should return something like this:
// A + 1 + X
// A + 1 + Y
// A + 1 + Z
// A + 2 + X
// ...
// B + 4 + Z
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

		if columns[0].OrdinateIndex == len(ssp) || completed {
			break
		}
	}

	return rows, nil
}

func gotoNextColumn(points *[]string, row *string, abscissaIndex int, column *Column, columns []Column, completed *bool) error {
	ssp, ok := column.Points.([]string)
	if !ok {
		return NewErrUnsopportedPointType()
	}

	if abscissaIndex+1 < len(columns) { // Until the last column is reached

		*row += ssp[column.OrdinateIndex]

		// Move forward
		abscissaIndex++
		column = &columns[abscissaIndex]
		err := gotoNextColumn(points, row, abscissaIndex, column, columns, completed)
		if err != nil {
			return err
		}

	} else { // When the last column is reached

		for _, point := range ssp {
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
	ssp, ok := column.Points.([]string)
	if !ok {
		return NewErrUnsopportedPointType()
	}

	if column.OrdinateIndex+1 < len(ssp) {
		column.OrdinateIndex++
	} else {
		column.OrdinateIndex = 0
		abscissaIndex--

		if abscissaIndex >= 0 {

			err := scrollDownPrevColumnPoint(&columns[abscissaIndex], columns, abscissaIndex, completed)
			if err != nil {
				return err
			}
		} else {
			*completed = true
		}
	}
	return nil
}
