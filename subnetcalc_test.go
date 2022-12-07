package subnetcalc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parse(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		_, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")
	})

	t.Run("Invalid", func(t *testing.T) {
		_, err := Parse("10.0.0.0/40")
		assert.ErrorIs(t, err, ErrCouldNotParse)
	})
}

func Test_SubnetInfo(t *testing.T) {
	s, err := Parse("10.0.0.0/16")
	assert.NoError(t, err, "parse should return no error")

	assert.Equal(t, "10.0.0.1", s.FirstIP())
	assert.Equal(t, "10.0.255.254", s.LastIP())
	assert.Equal(t, 16, s.Size())
}

func Test_Children(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")
		assert.Equal(t, false, s.HasChildReservations())
	})

	t.Run("Non-zero", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		_, err = s.AddReservation("10.0.0.0/24", "test")
		assert.NoError(t, err)

		_, err = s.AddReservation("10.0.1.0/24", "test")
		assert.NoError(t, err)

		assert.Equal(t, true, s.HasChildReservations())
	})
}

func Test_ReservationFail(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		_, err = s.AddReservation("192.168.0.0/24", "test")
		assert.ErrorIs(t, err, ErrDidNotFindSubnet)
	})

	t.Run("Already reserved", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		_, err = s.AddReservation("10.0.0.0/24", "test")
		assert.NoError(t, err)

		_, err = s.AddReservation("10.0.0.0/24", "test-2")
		assert.ErrorIs(t, ErrAlreadyReserved, err)
	})
}

func Test_UnReserve(t *testing.T) {
	t.Run("Not reserved", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		f, err := s.FindFree(17)
		assert.NoError(t, err)

		err = f.UnReserve()
		assert.ErrorIs(t, err, ErrNotReserved)
		assert.False(t, f.HasChildReservations())
		assert.Equal(t, f.subReservations, 0)
	})

	t.Run("Can reserve and update parent", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		f, err := s.FindFree(17)
		assert.NoError(t, err)

		err = f.Reserve("test")
		assert.NoError(t, err)
		assert.True(t, s.HasChildReservations())
		assert.Equal(t, s.subReservations, 1)
	})

	t.Run("Can unreserve and update parent", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		f, err := s.FindFree(17)
		assert.NoError(t, err)

		err = f.Reserve("test")
		assert.NoError(t, err)
		assert.True(t, s.HasChildReservations())
		assert.Equal(t, s.subReservations, 1)

		err = f.UnReserve()
		assert.NoError(t, err)
		assert.False(t, s.HasChildReservations())
		assert.Equal(t, s.subReservations, 0)

	})
}

func Test_FindFree(t *testing.T) {

	t.Run("/24 + /24", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/24", "reserved24-1")
		s.AddReservation("10.0.1.0/24", "reserved24-2")

		free, err := s.FindFree(24)
		assert.NoError(t, err)
		assert.Equal(t, "10.0.2.0/24", free.CIDR())
	})

	t.Run("/24 + /28", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/24", "reserved24-1")
		s.AddReservation("10.0.1.0/24", "reserved24-2")
		s.AddReservation("10.0.2.0/28", "reserved28")

		free, err := s.FindFree(28)
		assert.NoError(t, err)
		assert.Equal(t, "10.0.2.16/28", free.CIDR())
	})

	t.Run("/24 + /28 + /24", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/24", "reserved24-1")
		s.AddReservation("10.0.1.0/24", "reserved24-2")
		s.AddReservation("10.0.2.0/28", "reserved28-1")
		s.AddReservation("10.0.2.16/28", "reserved28-2")

		free, err := s.FindFree(24)
		assert.NoError(t, err)
		assert.Equal(t, "10.0.3.0/24", free.CIDR())
	})

	t.Run("/24 + /28 + /24 + /28", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/24", "reserved24-1")
		s.AddReservation("10.0.1.0/24", "reserved24-2")
		s.AddReservation("10.0.2.0/28", "reserved28-1")
		s.AddReservation("10.0.2.16/28", "reserved28-2")

		// Use reserve instead of addReservation
		free, err := s.FindFree(24)
		assert.NoError(t, err)
		assert.Equal(t, "10.0.3.0/24", free.CIDR())
		assert.NoError(t, free.Reserve("reserved24-3"))

		free, err = s.FindFree(28)
		assert.NoError(t, err)
		assert.Equal(t, "10.0.2.32/28", free.CIDR())
	})
}

func Test_FindFreeErrors(t *testing.T) {

	t.Run("propagated up", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/17", "some subnet")

		free, err := s.FindFree(16)
		assert.Nil(t, free)
		assert.Error(t, err)
	})

	t.Run("siblings", func(t *testing.T) {
		s, err := Parse("10.0.0.0/16")
		assert.NoError(t, err, "parse should return no error")

		s.AddReservation("10.0.0.0/17", "low")
		s.AddReservation("10.0.128.0/17", "high")

		free, err := s.FindFree(17)
		assert.Nil(t, free)
		assert.Error(t, err, ErrDidNotFindSubnet)
	})
}
