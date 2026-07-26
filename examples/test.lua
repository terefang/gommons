#!lua

arg['b'] = {1,2,3}
arg['b']['x'] = {1,2,3}
print('Hello World!',#arg,stringify(arg,1))
print('>>> ipairs')
for k,v in ipairs(arg) do
    print(k,v)
end
print('>>> apairs')
for k,v in apairs(arg) do
    print(k,v)
end
print('>>> pairs')
for k,v in pairs(arg) do
    print(k,v)
end
print('>>> _G')
for k,v in pairs(_G) do
    print(k,v)
end
print('>>> golib')
for k,v in pairs(golib) do
    print(k,v)
end