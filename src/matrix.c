#include "matrix.h"
#include <stdio.h>

Matrix createMatrix(void)
{
    Matrix board;
    clearFretBoard(&board);
    return board;
}

void clearFretBoard(Matrix *board)
{
    for (int i = 0; i < ROWS; i++)
        for (int j = 0; j < COLS; j++)
            board->matrix[i][j] = '-';
}

void updateFretBoard(Matrix *matrix, int row, int col, char character)
{
    if (row >= 0 && row < ROWS && col >= 0 && col < COLS)
        matrix->matrix[row][col] = character;
}

void matrixPrinter(const Matrix *matrix)
{
    printf("\n");
    for (int i = 0; i < ROWS; i++) {
        printf("|");
        for (int j = 0; j < COLS; j++)
            printf("%c", matrix->matrix[i][j]);
        printf("|\n");
    }
    printf("\n");
}
