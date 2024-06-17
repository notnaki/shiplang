fn m_func(n:number) {
    let prime = []number{};

    for (let i=0; i<n; i++;){
        prime = prime.append(1);
    }

    let p = 2;

    while(p*p <= n){
        if(prime[p] == 1) {
            for (let i = p*p; i < n; i = i + p;){
                prime[i] = 0;
            }
        }
        p++;
    }

    let primes = []number{};

    for (let p = 2; p < n; p++;){
        if (prime[p] == 1){
            primes = primes.append(p);
        }
    }

    return primes;
}


show(m_func(1000000));




