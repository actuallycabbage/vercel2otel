# api/
Vercel has a somewhat unusual way of running Go applications. All routes need 
to be sitting in the root api/ directory with a single `Handler` function 
in there.
