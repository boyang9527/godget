package collection

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TSDCache", func() {
	var (
		cache    *TSDCache
		capacity int
		err      interface{}
	)

	JustBeforeEach(func() {
		defer func() {
			err = recover()
		}()
		cache = NewTSDCache(capacity)
	})

	Describe("NewTSDCache", func() {
		Context("when creating TSDCache with invalid capacity", func() {
			BeforeEach(func() {
				capacity = -1

			})
			It("panics", func() {
				Expect(err).To(Equal("invalid TSDCache capacity"))
			})
		})
		Context("when creating TSDCache with valid capacity", func() {
			BeforeEach(func() {
				capacity = 10

			})
			It("returns the TSDCache", func() {
				Expect(err).To(BeNil())
				Expect(cache).NotTo(BeNil())
			})
		})
	})

	Describe("Put", func() {
		Context("when cache capacity is 1", func() {
			BeforeEach(func() {
				capacity = 1
			})
			It("only caches the latest data", func() {
				cache.Put(TestTSD{10})
				Expect(cache.String()).To(Equal("[{10}]"))
				cache.Put(TestTSD{20})
				Expect(cache.String()).To(Equal("[{20}]"))
				cache.Put(TestTSD{15})
				Expect(cache.String()).To(Equal("[{20}]"))
				cache.Put(TestTSD{30})
				Expect(cache.String()).To(Equal("[{30}]"))
			})
		})
		Context("when data put to cache do not execeed the capacity", func() {
			BeforeEach(func() {
				capacity = 5
			})
			It("cache all data in ascending order", func() {
				cache.Put(TestTSD{20})
				cache.Put(TestTSD{10})
				cache.Put(TestTSD{40})
				cache.Put(TestTSD{50})
				cache.Put(TestTSD{30})
				Expect(cache.String()).To(Equal("[{10} {20} {30} {40} {50}]"))
			})
		})
		Context("when data put to cache execeed the capacity", func() {
			BeforeEach(func() {
				capacity = 3
			})
			It("caches latest data in ascending order", func() {
				cache.Put(TestTSD{20})
				Expect(cache.String()).To(Equal("[{20}]"))
				cache.Put(TestTSD{10})
				Expect(cache.String()).To(Equal("[{10} {20}]"))
				cache.Put(TestTSD{40})
				Expect(cache.String()).To(Equal("[{10} {20} {40}]"))
				cache.Put(TestTSD{50})
				Expect(cache.String()).To(Equal("[{20} {40} {50}]"))
				cache.Put(TestTSD{30})
				Expect(cache.String()).To(Equal("[{30} {40} {50}]"))
				cache.Put(TestTSD{50})
				Expect(cache.String()).To(Equal("[{40} {50} {50}]"))
			})
		})

	})

	Describe("Query", func() {
		Context("when cache is empty", func() {
			It("return empty results", func() {
				result, ok := cache.Query(0, time.Now().UnixNano())
				Expect(ok).To(BeTrue())
				Expect(result).To(BeEmpty())
			})
		})
		Context("when data put to cache do not execeeds the capacity", func() {
			BeforeEach(func() {
				capacity = 5
			})
			It("returns the data in [start, end)", func() {
				cache.Put(TestTSD{20})
				result, ok := cache.Query(10, 40)
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal([]TSD{TestTSD{20}}))

				cache.Put(TestTSD{10})
				result, ok = cache.Query(10, 40)
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal([]TSD{TestTSD{10}, TestTSD{20}}))

				cache.Put(TestTSD{40})
				result, ok = cache.Query(10, 40)
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal([]TSD{TestTSD{10}, TestTSD{20}}))

				cache.Put(TestTSD{30})
				result, ok = cache.Query(10, 40)
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal([]TSD{TestTSD{10}, TestTSD{20}, TestTSD{30}}))

				cache.Put(TestTSD{50})
				result, ok = cache.Query(10, 40)
				Expect(ok).To(BeTrue())
				Expect(result).To(Equal([]TSD{TestTSD{10}, TestTSD{20}, TestTSD{30}}))
			})
		})

		Context("when data put to cache execeed the capacity", func() {
			BeforeEach(func() {
				capacity = 3
			})

			Context("when all queried data are guarenteed  in cache", func() {
				It("returns the data in [start, end)", func() {
					cache.Put(TestTSD{20})
					cache.Put(TestTSD{10})
					cache.Put(TestTSD{40})
					cache.Put(TestTSD{30})

					result, ok := cache.Query(30, 50)
					Expect(ok).To(BeTrue())
					Expect(result).To(Equal([]TSD{TestTSD{30}, TestTSD{40}}))

					cache.Put(TestTSD{50})
					result, ok = cache.Query(35, 50)
					Expect(ok).To(BeTrue())
					Expect(result).To(Equal([]TSD{TestTSD{40}}))
				})

			})
			Context("when queried data is possibly not in cache", func() {
				It("returns false", func() {
					cache.Put(TestTSD{20})
					cache.Put(TestTSD{10})
					cache.Put(TestTSD{40})
					cache.Put(TestTSD{30})

					_, ok := cache.Query(10, 50)
					Expect(ok).To(BeFalse())

					cache.Put(TestTSD{50})
					_, ok = cache.Query(30, 50)
					Expect(ok).To(BeFalse())
				})

			})

		})

	})
})
