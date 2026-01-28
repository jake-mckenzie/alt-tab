#ifndef MATRIX_H
#define MATRIX_H

#define ROWS 6
#define COLS 28

typedef struct {
    char matrix[ROWS][COLS];
} Matrix;

Matrix createMatrix(void);
void clearFretBoard(Matrix *board);
void updateFretBoard(Matrix *matrix, int row, int col, char character);
void matrixPrinter(const Matrix *matrix);

#endif
