import time
import multiprocessing
from jose import jwt

SECRET_KEY = "I think maybe python is over"
ROLE = "sub"

def generate_and_verify_jwt():
    sub = "1234567890"
    claims = {
        "sub": sub,
        "exp": int(time.time()) + 600,  # 10 minutes
        "role": ROLE
    }
    
    token = jwt.encode(claims, SECRET_KEY, algorithm="HS256")
    jwt.decode(token, SECRET_KEY, algorithms=["HS256"])

def run_test(count):
    for _ in range(count):
        generate_and_verify_jwt()

def run_parallel_test(count, parallelism):
    pool = multiprocessing.Pool(processes=parallelism)
    chunk_size = count // parallelism
    results = [pool.apply_async(run_test, (chunk_size,)) for _ in range(parallelism)]
    [result.get() for result in results]
    pool.close()
    pool.join()

if __name__ == "__main__":
    count = 100_000
    parallelisms = [1, 2, 4, 8, 16, 32, 64, 128]

    for p in parallelisms:
        start_time = time.time()
        if p == 1:
            run_test(count)
        else:
            run_parallel_test(count, p)
        duration = time.time() - start_time
        print(f"Parallelism: {p}, Time: {duration:.2f} seconds")
