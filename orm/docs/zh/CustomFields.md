## Custom Fields

	TypeBooleanField = 1 << iota

	// string
	TypeCharField

	// string
	TypeTextField

	// time.Time
	TypeDateField
	// time.Time
	TypeDateTimeField

	// int16
	TypeSmallIntegerField
	// int32
	TypeIntegerField
	// int64
	TypeBigIntegerField
	// uint16
	TypePositiveSmallIntegerField
	// uint32
	TypePositiveIntegerField
	// uint64
	TypePositiveBigIntegerField

	// float64
	TypeFloatField
	// float64
	TypeDecimalField

	RelForeignKey
	RelOneToOne
	RelManyToMany
	RelReverseOne
	RelReverseMany