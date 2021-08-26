(
  (
    "not of nil is t"
    (import ("core/not")
      (not ()))
    t
  )
  (
    "not of t is nil"
    (import ("core/not")
      (not (quote t)))
    ()
  )
  (
    "not of number is nil"
    (import ("core/not")
      (not 1))
    ()
  )
  (
    "not with no args"
    (import ("core/not")
      (label actual (try (not))
        (cond
	  (car actual) actual
	  (quote t) (quote t))))
    t
  )
  (
    "not with many args"
    (import ("core/not")
      (label actual (try (not 1 2))
        (cond
          (car actual) actual
          (quote t) (quote t))))
    t
  )
)
