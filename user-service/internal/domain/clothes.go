package domain

type Clothes struct {
	ID      int64
	Name    string
	Price   float64
	ColorID *int64
	Images  []ClothesImage
}

type ClothesImage struct {
	ID        int64
	ClothesID int64
	ImageURL  string
}
