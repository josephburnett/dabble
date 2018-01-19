(
  (
    "and t is t (does not work)"
    (car (try
    (import ("core/and")
      (and (quote t)))))
    ()
  )
  (
    "and single true"
    (import ("core/and")
      (and 1))
    t
  )
  (
    "and single false"
    (import ("core/and")
      (and ()))
    ()
  ) 
  (
    "and multiple true"
    (import ("core/and")
      (and 1 2 3 4))
    t
  )
  (
    "and multiple false"
    (import ("core/and")
      (and 1 2 3 ()))
    ()
  )
)