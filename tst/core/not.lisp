(if (eq () (not t))
    (if (eq t (not ()))
	t
      (error "(not () must return t"))
  (error "(not t) must return ()"))

  
