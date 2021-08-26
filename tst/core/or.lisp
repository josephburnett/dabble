(
  (
    "or t is t (does not work)"
    (car (try
    (import ("core/or")
      (or (quote t)))))
    ()
  )
  (
    "or nil is nil"
    (import ("core/or")
      (or ()))
    ()
  )
  (
    "or many true is t"
    (import ("core/or")
      (or 1 2 3 4))
    t
  )
  (
    "or many false is nil"
    (import ("core/or")
      (or () ()))
    ()
  )
  (
    "or one of many false is true"
    (import ("core/or")
      (or () 1 () 2))
    t
  )
)
