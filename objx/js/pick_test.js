// Test cases for pick.js
const { pick } = require('./pick');

// 辅助函数，用于清晰显示测试结果
function printTestResult(testName, actual, expected) {
  console.log(`测试: ${testName}`);
  console.log(`期望找到${expected}个项目，实际找到: ${actual.length}`);
  
  // 简单格式化输出结果，避免太复杂的嵌套结构
  const simplifyForDisplay = (item) => {
    if (Array.isArray(item)) {
      return `Array(${item.length})`;
    }
    if (item !== null && typeof item === 'object') {
      const keys = Object.keys(item);
      if (keys.includes('name')) {
        return `{name: "${item.name}", ...}`;
      }
      if (keys.includes('id')) {
        return `{id: ${item.id}, ...}`;
      }
      return `{${keys.join(', ')}}`;
    }
    return item;
  };
  
  console.log("匹配结果:", actual.map(simplifyForDisplay));
  console.log("\n");
}

// 测试1：基本选择器测试
function testBasicSelectors() {
  console.log("==== 基本选择器测试 ====");
  
  const data = {
    users: [
      { id: 1, name: "张三", age: 30, roles: ["admin", "user"] },
      { id: 2, name: "李四", age: 25, roles: ["user"] },
      { id: 3, name: "王五", age: 35, roles: ["manager", "user"] }
    ],
    settings: {
      theme: "dark",
      notifications: true
    }
  };
  
  // 测试 1.1: 选择名字为"张三"的用户
  const test1_1 = pick(data, "users [name=张三]");
  printTestResult("1.1 选择名字为'张三'的用户", test1_1, 1);
  
  // 测试 1.2: 选择年龄为30的用户
  const test1_2 = pick(data, "users [age=30]");
  printTestResult("1.2 选择年龄为30的用户", test1_2, 1);
  
  // 测试 1.3: 模糊匹配名字包含"三"的用户
  const test1_3 = pick(data, "users [name*=三]");
  printTestResult("1.3 模糊匹配名字包含'三'的用户", test1_3, 1);
}

// 测试2：复杂对象和嵌套结构测试
function testComplexObjects() {
  console.log("==== 复杂对象和嵌套结构测试 ====");
  
  const data = {
    departments: [
      {
        name: "技术部",
        employees: [
          { id: 101, name: "张三", position: "开发工程师", skills: ["JavaScript", "React", "Node.js"] },
          { id: 102, name: "李四", position: "测试工程师", skills: ["Python", "TestNG"] }
        ]
      },
      {
        name: "市场部",
        employees: [
          { id: 201, name: "王五", position: "市场经理", skills: ["市场分析", "客户关系"] },
          { id: 202, name: "赵六", position: "销售代表", skills: ["谈判", "销售"] }
        ]
      }
    ]
  };
  
  // 测试 2.1: 选择技术部的所有员工
  const test2_1 = pick(data, "departments [name=技术部] employees");
  printTestResult("2.1 选择技术部的所有员工", test2_1, 2);
  
  // 测试 2.2: 选择具有JavaScript技能的员工
  const test2_2 = pick(data, "departments employees [skills*=JavaScript]");
  printTestResult("2.2 选择具有JavaScript技能的员工", test2_2, 1);
  
  // 测试 2.3: 选择不是销售代表的员工
  const test2_3 = pick(data, "departments employees [position!=销售代表]");
  printTestResult("2.3 选择不是销售代表的员工", test2_3, 3);
}

// 测试3：数组和原始类型测试
function testArraysAndPrimitives() {
  console.log("==== 数组和原始类型测试 ====");
  
  const data = [
    "apple",
    "banana",
    "orange",
    { type: "fruit", name: "grape", color: "purple" },
    { type: "vegetable", name: "carrot", color: "orange" }
  ];
  
  // 测试 3.1: 选择所有水果
  const test3_1 = pick(data, "[type=fruit]");
  printTestResult("3.1 选择所有水果", test3_1, 1);
  
  // 测试 3.2: 选择橙色的项目
  const test3_2 = pick(data, "[color=orange]");
  printTestResult("3.2 选择橙色的项目", test3_2, 1);
}

// 测试4：边缘情况测试
function testEdgeCases() {
  console.log("==== 边缘情况测试 ====");
  
  // 测试 4.1: 空数据
  const test4_1 = pick({}, "users");
  printTestResult("4.1 空数据", test4_1, 0);
  
  // 测试 4.2: 空选择器
  const data = { users: [{ name: "张三" }] };
  const test4_2 = pick(data, "");
  printTestResult("4.2 空选择器", test4_2, 1);
  
  // 测试 4.3: null数据
  const test4_3 = pick(null, "users");
  printTestResult("4.3 null数据", test4_3, 0);
}

// 运行所有测试
function runAllTests() {
  testBasicSelectors();
  testComplexObjects();
  testArraysAndPrimitives();
  testEdgeCases();
  
  console.log("所有测试完成！");
}

runAllTests(); 