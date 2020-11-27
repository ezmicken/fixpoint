// Package fixpoint implements fixed-point arithmetic and vector operations. It
// has been inspired on and is partially copied from the github.com/go-gl/mathgl
// package.
package fixpoint

// Useful link:
// https://spin.atomicobject.com/2012/03/15/simple-fixed-point-math/

// Q16 is a Q7.16 fixed point integer type that has 16 bits of precision to the
// right of the fixed point. It is designed to be used as a more efficient
// replacement for unit vectors with some extra room to avoid overflow.
type Q16 struct {
	N int32
}

// Q16FromFloat converts a float32 to the same number in fixed point format.
// Inverse of .Float().
func Q16FromFloat(x float32) Q16 {
	return Q16{int32(x * (1 << 16))}
}

// Q16FromInt32 returns a fixed point integer with all decimals set to zero.
func Q16FromInt32(x int32) Q16 {
	return Q16{x << 16}
}

// Float returns the floating point version of this fixed point number. Inverse
// of Q16FromFloat.
func (q Q16) Float() float32 {
	return float32(q.N) / (1 << 16)
}

// Int32Scaled returns the underlying fixed point number multiplied by scale.
func (q Q16) Int32Scaled(scale int32) int32 {
	return q.N / (1 << 16 / scale)
}

// Add returns the argument plus this number.
func (q1 Q16) Add(q2 Q16) Q16 {
	return Q16{q1.N + q2.N}
}

// Sub returns the argument minus this number.
func (q1 Q16) Sub(q2 Q16) Q16 {
	return Q16{q1.N - q2.N}
}

// Neg returns the inverse of this number.
func (q1 Q16) Neg() Q16 {
	return Q16{-q1.N}
}

// Mul returns this number multiplied by the argument.
func (q1 Q16) Mul(q2 Q16) Q16 {
	return Q16{int32((int64(q1.N) * int64(q2.N)) >> 16)}
}

// Div returns this number divided by the argument.
func (q1 Q16) Div(q2 Q16) Q16 {
	return Q16{int32((int64(q1.N) << 16) / int64(q2.N))}
}

// Vec3Q16 is a 3-dimensional vector with Q16 fixed point elements.
type Vec3Q16 struct {
	X Q16
	Y Q16
	Z Q16
}

// Vec3Q16FromFloat returns the fixed-point vector of the given 3 floats.
func Vec3Q16FromFloat(x, y, z float32) Vec3Q16 {
	return Vec3Q16{Q16FromFloat(x), Q16FromFloat(y), Q16FromFloat(z)}
}

// Add returns this vector added to the argument.
func (v1 Vec3Q16) Add(v2 Vec3Q16) Vec3Q16 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q16{v1.X.Add(v2.X), v1.Y.Add(v2.Y), v1.Z.Add(v2.Z)}
}

// Mul returns this vector multiplied by the argument.
func (v1 Vec3Q16) Mul(c Q16) Vec3Q16 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q16{v1.X.Mul(c), v1.Y.Mul(c), v1.Z.Mul(c)}
}

// Dot returns the dot product between this vector and the argument.
func (v1 Vec3Q16) Dot(v2 Vec3Q16) Q16 {
	// Copied from go-gl/mathgl and modified.
	return v1.X.Mul(v2.X).Add(v1.Y.Mul(v2.Y)).Add(v1.Z.Mul(v2.Z))
}

func (v1 Vec3Q16) Sub(v2 Vec3Q16) Vec3Q16 {
  return Vec3Q16{v1.X.Sub(v2.X), v1.Y.Sub(v2.Y), v1.Z.Sub(v2.Z)}
}

// Cross returns the cross product between this vector and the argument.
func (v1 Vec3Q16) Cross(v2 Vec3Q16) Vec3Q16 {
	// Copied from go-gl/mathgl and modified.
	return Vec3Q16{v1.Y.Mul(v2.Z).Sub(v1.Z.Mul(v2.Y)), v1.Z.Mul(v2.X).Sub(v1.X.Mul(v2.Z)), v1.X.Mul(v2.Y).Sub(v1.Y.Mul(v2.X))}
}

// QuatQ16 is a quaternion with Q16 fixed point elements.
type QuatQ16 struct {
	W Q16
	V Vec3Q16
}

// QuatIdent returns the identity quaternion.
func QuatIdent() QuatQ16 {
	return QuatQ16{Q16FromInt32(1), Vec3Q16{}}
}

// X returns the X part of this quaternion.
func (q QuatQ16) X() Q16 {
	return q.V.X
}

// Y returns the Y part of this quaternion.
func (q QuatQ16) Y() Q16 {
	return q.V.Y
}

// Z returns the Z part of this quaternion.
func (q QuatQ16) Z() Q16 {
	return q.V.Z
}

// Mul returns this quaternion multiplied by the argument.
func (q1 QuatQ16) Mul(q2 QuatQ16) QuatQ16 {
	// Copied from go-gl/mathgl and modified.
	return QuatQ16{q1.W.Mul(q2.W).Sub(q1.V.Dot(q2.V)), q1.V.Cross(q2.V).Add(q2.V.Mul(q1.W)).Add(q1.V.Mul(q2.W))}
}

// Rotate returns the vector from the argument rotated by the rotation this
// quaternion represents.
func (q1 QuatQ16) Rotate(v Vec3Q16) Vec3Q16 {
	// Copied from go-gl/mathgl and modified.
	cross := q1.V.Cross(v)
	// v + 2q_w * (q_v x v) + 2q_v x (q_v x v)
	return v.Add(cross.Mul(Q16FromInt32(2).Mul(q1.W))).Add(q1.V.Mul(Q16FromInt32(2)).Cross(cross))
}
