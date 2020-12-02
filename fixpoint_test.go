package fixpoint

import (
  "testing"
  "math"

  "github.com/go-gl/mathgl/mgl32"
  "github.com/stretchr/testify/assert"
)

func TestQ16(t *testing.T) {
  two := Q16FromFloat(2)
  for _, f := range []float32{0.25, 1, 10, 0.125} {
    q := Q16FromFloat(f)
    assert.Equal(t, f, q.Float(), "float32 roundtrip failed")
    assert.Equal(t, Q16FromFloat(f*2), q.Mul(two), "multiply by 2")
    square := q.Mul(q)
    assert.Equal(t, Q16FromFloat(f*f), square, "square")
    assert.Equal(t, Q16FromFloat(f), square.Div(q), "div")
  }
}

func TestRotation(t *testing.T) {
  // Create vector to rotate.
  vec1 := mgl32.Vec3{0, 0.832, 0.554}
  vec2 := Vec3Q16FromFloat(0, 0.832, 0.554)

  // Create rotation quaternion.
  inc1 := mgl32.Quat{0, mgl32.Vec3{0.07, 0, 0}}
  inc1.W = 1.0 - 0.5*(inc1.X()*inc1.X()+inc1.Y()*inc1.Y()+inc1.Z()*inc1.Z())

  inc2 := QuatQ16{Q16FromInt32(0), Vec3Q16FromFloat(0.07, 0, 0)}
  inc2.W = Q16FromInt32(1).Sub(Q16FromFloat(0.5).Mul(inc2.X().Mul(inc2.X()).Add(inc2.Y().Mul(inc2.Y())).Add(inc2.Z().Mul(inc2.Z()))))

  rotation1 := mgl32.QuatIdent()
  rotation2 := QuatIdent()
  for i := 0; i < 10; i++ {
    // Incrementally add parts of the rotation.
    rotation1 = rotation1.Mul(inc1)
    rotation2 = rotation2.Mul(inc2)

    // Get the vector that belongs to this rotation.
    vec1 := rotation1.Rotate(vec1)
    vec2 := rotation2.Rotate(vec2)
    diffX := vec2.X.Float() - vec1.X()
    diffY := vec2.Y.Float() - vec1.Y()
    diffZ := vec2.Z.Float() - vec1.Z()
    if diffX > 0.001 || diffY > 0.001 || diffZ > 0.001 {
      t.Errorf("difference too big:\nX1: %.8f\nX2: %.8f\nY1: %.8f\nY2: %.8f\nZ1: %.8f\nZ2: %.8f",
        vec1.X(), vec2.X.Float(),
        vec1.Y(), vec2.Y.Float(),
        vec1.Z(), vec2.Z.Float())
    }
  }
}

func TestInvSqrt(t *testing.T) {
  numbers := []Q16{
    Q16FromFloat(3.948),
    Q16FromFloat(11.039),
    Q16FromFloat(13.481),
    Q16FromFloat(12.677),
    Q16FromFloat(19.708),
    Q16FromFloat(9.198),
    Q16FromFloat(4.204),
    Q16FromFloat(12.566),
    Q16FromFloat(14.606),
    Q16FromFloat(19.896),
    Q16FromFloat(14.038),
    Q16FromFloat(19.073),
    Q16FromFloat(8.127),
    Q16FromFloat(15.579),
    Q16FromFloat(4.623),
    Q16FromFloat(18.316),
    Q16FromFloat(16.349),
    Q16FromFloat(14.642),
    Q16FromFloat(3.707),
    Q16FromFloat(19.464),
    Q16FromFloat(9.104),
    Q16FromFloat(8.587),
    Q16FromFloat(14.378),
    Q16FromFloat(9.098),
    Q16FromFloat(6.048),
    Q16FromFloat(10.968),
    Q16FromFloat(8.103),
    Q16FromFloat(2.142),
    Q16FromFloat(2.967),
    Q16FromFloat(7.559),
    Q16FromFloat(7.851),
    Q16FromFloat(10.291),
    Q16FromFloat(6.629),
    Q16FromFloat(18.558),
    Q16FromFloat(15.93),
    Q16FromFloat(17.584),
    Q16FromFloat(14.214),
    Q16FromFloat(3.036),
    Q16FromFloat(9.574),
    Q16FromFloat(10.9),
    Q16FromFloat(2.305),
    Q16FromFloat(14.024),
    Q16FromFloat(15.022),
    Q16FromFloat(9.1),
    Q16FromFloat(14.996),
    Q16FromFloat(7.626),
    Q16FromFloat(6.321),
    Q16FromFloat(13.596),
    Q16FromFloat(4.792),
    Q16FromFloat(6.587),
    Q16FromFloat(7.631),
    Q16FromFloat(12.021),
    Q16FromFloat(16.673),
    Q16FromFloat(16.424),
    Q16FromFloat(2.883),
    Q16FromFloat(19.8),
    Q16FromFloat(16.912),
    Q16FromFloat(8.707),
    Q16FromFloat(15.343),
    Q16FromFloat(4.478),
    Q16FromFloat(14.4),
    Q16FromFloat(10.646),
    Q16FromFloat(12.422),
    Q16FromFloat(12.058),
    Q16FromFloat(7.388),
    Q16FromFloat(16.013),
    Q16FromFloat(11.235),
    Q16FromFloat(11.515),
    Q16FromFloat(16.884),
    Q16FromFloat(16.384),
    Q16FromFloat(16.738),
    Q16FromFloat(10.365),
    Q16FromFloat(11.491),
    Q16FromFloat(3.038),
    Q16FromFloat(15.855),
    Q16FromFloat(11.239),
    Q16FromFloat(5.458),
    Q16FromFloat(19.413),
    Q16FromFloat(10.805),
    Q16FromFloat(19.227),
    Q16FromFloat(11.748),
    Q16FromFloat(10.569),
    Q16FromFloat(12.489),
    Q16FromFloat(15.134),
    Q16FromFloat(10.83),
    Q16FromFloat(4.953),
    Q16FromFloat(14.684),
    Q16FromFloat(7.554),
    Q16FromFloat(14.624),
    Q16FromFloat(16.185),
    Q16FromFloat(8.856),
    Q16FromFloat(8.575),
    Q16FromFloat(2.58),
    Q16FromFloat(17.272),
    Q16FromFloat(15.103),
    Q16FromFloat(8.976),
    Q16FromFloat(6.556),
    Q16FromFloat(19.019),
    Q16FromFloat(12.541),
    Q16FromFloat(4.554),
  }

  for i := 0; i < len(numbers); i++ {
    floatAns := float32(float64(1)/math.Sqrt(float64(numbers[i].Float())))
    fixedAns := numbers[i].InvSqrt().Float()
    diff := mgl32.Abs(floatAns - fixedAns)
    if diff > 0.001 {
      t.Errorf("difference too big: %v: %v vs %v (%v)", numbers[i], floatAns, fixedAns, diff)
    }
  }
}
