(cons (cons (quote and)
            (cons (macro (x xs)
                         (cond
                           x (cond
                               (eq () xs) (quote t)
                               (quote t) (and xs))
                           (quote t) ())) ())) ())
