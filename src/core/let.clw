(cons (cons (quote let)
            (cons (macro (binding form xs)
                         (cond
                           (eq () (atom xs)) (error)
                           (eq () binding) (error)
                           (quote t) (label (car (car binding))
                                            (car (cdr (car binding)))
                                            (cond
                                              (eq () (cdr binding)) form
                                              (quote t) (recur (cdr binding) form))))) ())) ())
