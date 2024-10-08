package transform

type transform struct{}

func NewTransform() *transform {
	return &transform{}
}

func (t *transform) TransformBlockToBlockOracle(blocks []any) (any, error) {
	return nil, nil
}
